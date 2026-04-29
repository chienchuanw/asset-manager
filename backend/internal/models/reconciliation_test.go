package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestReconcileBankAccountInput_Validate(t *testing.T) {
	t.Run("有效輸入通過驗證", func(t *testing.T) {
		in := &ReconcileBankAccountInput{NewBalance: 1000, Date: time.Now()}
		assert.NoError(t, in.Validate())
	})
	t.Run("負數餘額被拒絕", func(t *testing.T) {
		in := &ReconcileBankAccountInput{NewBalance: -1, Date: time.Now()}
		assert.Error(t, in.Validate())
	})
	t.Run("零值日期被拒絕", func(t *testing.T) {
		in := &ReconcileBankAccountInput{NewBalance: 0, Date: time.Time{}}
		assert.Error(t, in.Validate())
	})
}

func TestReconcileCreditCardInput_Validate(t *testing.T) {
	t.Run("至少需要提供一個新值", func(t *testing.T) {
		in := &ReconcileCreditCardInput{Date: time.Now()}
		assert.Error(t, in.Validate())
	})
	t.Run("僅提供 used_credit 是有效的", func(t *testing.T) {
		uc := 500.0
		in := &ReconcileCreditCardInput{NewUsedCredit: &uc, Date: time.Now()}
		assert.NoError(t, in.Validate())
	})
	t.Run("僅提供 credit_limit 是有效的", func(t *testing.T) {
		lim := 100000.0
		in := &ReconcileCreditCardInput{NewCreditLimit: &lim, Date: time.Now()}
		assert.NoError(t, in.Validate())
	})
	t.Run("負數 used_credit 被拒絕", func(t *testing.T) {
		neg := -1.0
		in := &ReconcileCreditCardInput{NewUsedCredit: &neg, Date: time.Now()}
		assert.Error(t, in.Validate())
	})
	t.Run("零或負數 credit_limit 被拒絕", func(t *testing.T) {
		zero := 0.0
		in := &ReconcileCreditCardInput{NewCreditLimit: &zero, Date: time.Now()}
		assert.Error(t, in.Validate())
	})
	t.Run("零值日期被拒絕", func(t *testing.T) {
		uc := 100.0
		in := &ReconcileCreditCardInput{NewUsedCredit: &uc, Date: time.Time{}}
		assert.Error(t, in.Validate())
	})
}

func TestReconcileBatchInput_Validate(t *testing.T) {
	t.Run("空 items 被拒絕", func(t *testing.T) {
		in := &ReconcileBatchInput{Date: time.Now(), Items: []ReconcileBatchItem{}}
		assert.Error(t, in.Validate())
	})
	t.Run("零值日期被拒絕", func(t *testing.T) {
		bal := 100.0
		in := &ReconcileBatchInput{
			Date: time.Time{},
			Items: []ReconcileBatchItem{
				{TargetType: SourceTypeBankAccount, TargetID: uuid.New(), NewBalance: &bal},
			},
		}
		assert.Error(t, in.Validate())
	})
	t.Run("bank_account item 缺 new_balance 被拒絕", func(t *testing.T) {
		in := &ReconcileBatchInput{
			Date: time.Now(),
			Items: []ReconcileBatchItem{
				{TargetType: SourceTypeBankAccount, TargetID: uuid.New()},
			},
		}
		assert.Error(t, in.Validate())
	})
	t.Run("credit_card item 沒有任何新值被拒絕", func(t *testing.T) {
		in := &ReconcileBatchInput{
			Date: time.Now(),
			Items: []ReconcileBatchItem{
				{TargetType: SourceTypeCreditCard, TargetID: uuid.New()},
			},
		}
		assert.Error(t, in.Validate())
	})
	t.Run("不支援的 target_type 被拒絕", func(t *testing.T) {
		bal := 100.0
		in := &ReconcileBatchInput{
			Date: time.Now(),
			Items: []ReconcileBatchItem{
				{TargetType: SourceTypeManual, TargetID: uuid.New(), NewBalance: &bal},
			},
		}
		assert.Error(t, in.Validate())
	})
	t.Run("有效混合 batch 通過驗證", func(t *testing.T) {
		bal := 100.0
		uc := 200.0
		in := &ReconcileBatchInput{
			Date: time.Now(),
			Items: []ReconcileBatchItem{
				{TargetType: SourceTypeBankAccount, TargetID: uuid.New(), NewBalance: &bal},
				{TargetType: SourceTypeCreditCard, TargetID: uuid.New(), NewUsedCredit: &uc},
			},
		}
		assert.NoError(t, in.Validate())
	})
}
