package mcpserver

import (
	"errors"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeHoldingSvc struct {
	all *service.HoldingServiceResult
	one *models.Holding
	err error
}

func (f fakeHoldingSvc) GetAllHoldings(models.HoldingFilters) (*service.HoldingServiceResult, error) {
	return f.all, f.err
}
func (f fakeHoldingSvc) GetHoldingBySymbol(string) (*models.Holding, error) { return f.one, f.err }

func TestHoldingsAdapter_List(t *testing.T) {
	a := NewHoldingsAdapter(fakeHoldingSvc{all: &service.HoldingServiceResult{
		Holdings: []*models.Holding{{Symbol: "2330", Name: "TSMC", Quantity: 10}},
	}})
	out, err := a.List()
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "2330", out[0].Symbol)
	assert.Equal(t, 10.0, out[0].Quantity)
}

func TestHoldingsAdapter_ListError(t *testing.T) {
	a := NewHoldingsAdapter(fakeHoldingSvc{err: errors.New("boom")})
	_, err := a.List()
	assert.ErrorContains(t, err, "boom")
}

type fakeTxSvc struct {
	list []*models.Transaction
	one  *models.Transaction
	err  error
}

func (f fakeTxSvc) ListTransactions(repository.TransactionFilters) ([]*models.Transaction, error) {
	return f.list, f.err
}
func (f fakeTxSvc) GetTransaction(uuid.UUID) (*models.Transaction, error) { return f.one, f.err }

func TestTransactionsAdapter_List(t *testing.T) {
	a := NewTransactionsAdapter(fakeTxSvc{list: []*models.Transaction{{Symbol: "AAPL"}}})
	out, err := a.List(repository.TransactionFilters{})
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "AAPL", out[0].Symbol)
}

func TestTransactionsAdapter_GetInvalidUUID(t *testing.T) {
	a := NewTransactionsAdapter(fakeTxSvc{})
	_, err := a.Get("not-a-uuid")
	assert.Error(t, err)
}

type fakeAllocSvc struct{ rows []models.AllocationByType }

func (f fakeAllocSvc) GetAllocationByType() ([]models.AllocationByType, error) { return f.rows, nil }

func TestAllocationAdapter_ByType(t *testing.T) {
	a := NewAllocationAdapter(fakeAllocSvc{rows: []models.AllocationByType{{}}})
	out, err := a.ByType()
	require.NoError(t, err)
	assert.Len(t, out, 1)
}
