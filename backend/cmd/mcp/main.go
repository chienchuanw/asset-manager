// Command mcp runs the in-process MCP server (Phase 1, read-only) over stdio.
// It reuses the same DB + service layer as cmd/api (no Redis cache path) and
// records every tool call to agent_audit_log.
package main

import (
	"context"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"github.com/chienchuanw/asset-manager/internal/audit"
	"github.com/chienchuanw/asset-manager/internal/client"
	"github.com/chienchuanw/asset-manager/internal/db"
	"github.com/chienchuanw/asset-manager/internal/mcpserver"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/chienchuanw/asset-manager/internal/service"
)

func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Printf("Warning: .env.local not found: %v", err)
	}

	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Repositories
	transactionRepo := repository.NewTransactionRepository(database)
	exchangeRateRepo := repository.NewExchangeRateRepository(database)
	realizedProfitRepo := repository.NewRealizedProfitRepository(database)
	dbx := sqlx.NewDb(database, "postgres")
	performanceSnapshotRepo := repository.NewPerformanceSnapshotRepository(dbx)

	// Price service (real API if keys present, else mock) — no Redis cache.
	var priceService service.PriceService
	finmind := os.Getenv("FINMIND_API_KEY")
	coingecko := os.Getenv("COINGECKO_API_KEY")
	alpha := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if finmind != "" && coingecko != "" && alpha != "" {
		priceService = service.NewRealPriceService(finmind, coingecko, alpha)
	} else {
		priceService = service.NewMockPriceService()
		log.Println("Warning: price API keys not set, using mock price service")
	}

	exchangeRateClient := client.NewExchangeRateAPIClient()
	exchangeRateService := service.NewExchangeRateService(exchangeRateRepo, exchangeRateClient, nil)
	fifoCalculator := service.NewFIFOCalculator(exchangeRateService)
	transactionService := service.NewTransactionService(transactionRepo, realizedProfitRepo, fifoCalculator, exchangeRateService)
	holdingService := service.NewHoldingService(transactionRepo, fifoCalculator, priceService, exchangeRateService)
	analyticsService := service.NewAnalyticsService(realizedProfitRepo)
	unrealizedAnalyticsService := service.NewUnrealizedAnalyticsService(holdingService)
	allocationService := service.NewAllocationService(holdingService)
	performanceTrendService := service.NewPerformanceTrendService(performanceSnapshotRepo, unrealizedAnalyticsService, analyticsService)

	reconcileService := service.NewHoldingReconciliationService(database, transactionRepo, fifoCalculator)

	auditRepo := audit.NewPostgresAuditRepository(database)
	logger := audit.NewLogger(auditRepo)
	failClosed := audit.NewFailClosed(auditRepo)

	holdingsAdapter := mcpserver.NewHoldingsAdapter(holdingService)
	txAdapter := mcpserver.NewTransactionsAdapter(transactionService)
	txWriteAdapter := mcpserver.NewTxWriteAdapter(transactionService)
	reconcileAdapter := mcpserver.NewReconcileAdapter(reconcileService)
	tools := []mcpserver.Tool{
		// read (Phase 1)
		mcpserver.NewGetHoldingsTool(holdingsAdapter),
		mcpserver.NewGetHoldingTool(holdingsAdapter),
		mcpserver.NewListTransactionsTool(txAdapter),
		mcpserver.NewGetTransactionTool(txAdapter),
		mcpserver.NewGetAnalyticsTool(mcpserver.NewAnalyticsAdapter(analyticsService)),
		mcpserver.NewGetAllocationTool(mcpserver.NewAllocationAdapter(allocationService)),
		mcpserver.NewGetPerformanceTrendTool(mcpserver.NewTrendAdapter(performanceTrendService)),
		// write (Phase 2, dry-run/confirm + fail-closed audit)
		mcpserver.NewCreateTransactionTool(txWriteAdapter, failClosed),
		mcpserver.NewUpdateTransactionTool(txWriteAdapter, failClosed),
		mcpserver.NewDeleteTransactionTool(txWriteAdapter, failClosed),
		mcpserver.NewReconcileHoldingTool(reconcileAdapter, failClosed),
	}

	dispatcher := mcpserver.NewDispatcher(tools, mcpserver.NewLoggerSink(logger))
	if err := mcpserver.Serve(context.Background(), dispatcher, tools); err != nil {
		log.Fatalf("mcp server: %v", err)
	}
}
