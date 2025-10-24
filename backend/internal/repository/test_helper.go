package repository

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// setupTestDB 設定測試資料庫連線
func setupTestDB() (*sql.DB, error) {
	// 從環境變數讀取測試資料庫設定
	host := getEnv("TEST_DB_HOST", "localhost")
	port := getEnv("TEST_DB_PORT", "5432")
	user := getEnv("TEST_DB_USER", "postgres")
	password := getEnv("TEST_DB_PASSWORD", "postgres")
	dbname := getEnv("TEST_DB_NAME", "asset_manager_test")

	// 建立連線字串
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// 連接資料庫
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 測試連線
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// getEnv 取得環境變數，如果不存在則使用預設值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// cleanupAssetSnapshots 清理測試資料庫中的資產快照資料
func cleanupAssetSnapshots(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM asset_snapshots")
	if err != nil {
		return fmt.Errorf("failed to cleanup asset_snapshots: %w", err)
	}
	return nil
}

// cleanupTransactions 清理測試資料庫中的交易記錄資料
func cleanupTransactions(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM transactions")
	if err != nil {
		return fmt.Errorf("failed to cleanup transactions: %w", err)
	}
	return nil
}

// cleanupAllTables 清理測試資料庫中的所有資料
func cleanupAllTables(db *sql.DB) error {
	// 按照外鍵依賴順序清理
	tables := []string{
		"asset_snapshots",
		"exchange_rates",
		"transactions",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("failed to cleanup table %s: %w", table, err)
		}
	}

	return nil
}
