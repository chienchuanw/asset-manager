# Cash Withdrawal Feature Implementation

## Overview

This document describes the implementation of the cash withdrawal feature, which allows users to record cash withdrawals from their bank accounts.

## Feature Description

- **Purpose**: Record cash withdrawals from bank accounts
- **Behavior**: 
  - Deducts the withdrawal amount from the bank account balance
  - Does not track the withdrawn cash
  - Uses existing `transfer_out` type with a new "提領" (withdrawal) category
  - Target type is set to `cash` with `target_id` as `null`

## Database Changes

### Migration Files

#### 1. `000024_add_withdrawal_category.up.sql`
- Adds "提領" (withdrawal) category for `transfer_out` type
- Category is marked as system category (`is_system = true`)
- Uses `ON CONFLICT DO NOTHING` to prevent duplicate insertion errors

#### 2. `000025_add_cash_source_type.up.sql`
- Updates `cash_flows.source_type` constraint to include `'cash'` type
- Updates column comments to document the new type
- Updates `target_type` comment to include cash as a valid target type

### Rollback Files

- `000024_add_withdrawal_category.down.sql`: Removes the withdrawal category
- `000025_add_cash_source_type.down.sql`: Removes cash from source type constraint

## Code Changes

### Backend

#### 1. Models (`backend/internal/models/cash_flow.go`)
```go
const (
    SourceTypeManual       SourceType = "manual"
    SourceTypeSubscription SourceType = "subscription"
    SourceTypeInstallment  SourceType = "installment"
    SourceTypeBankAccount  SourceType = "bank_account"
    SourceTypeCreditCard   SourceType = "credit_card"
    SourceTypeCash         SourceType = "cash"  // New
)
```

#### 2. Service (`backend/internal/service/cash_flow_service.go`)
- Updated `CreateCashFlow` to handle cash withdrawals
- When `target_type = "cash"`, no target account balance update is performed
- Validates that `target_id` must be `null` when `target_type = "cash"`

#### 3. Tests
- `backend/internal/api/cash_flow_handler_test.go`: Added `TestCreateCashFlow_CashWithdrawal`
- `backend/internal/service/cash_flow_balance_test.go`: Added `TestCashFlowService_CreateCashFlow_CashWithdrawal`
- `backend/internal/repository/category_repository_test.go`: Updated `ensureSystemCategories` to include withdrawal category

### Frontend

#### 1. Types (`frontend/src/types/cash-flow.ts`)
```typescript
export const SourceType = {
  MANUAL: "manual",
  SUBSCRIPTION: "subscription",
  INSTALLMENT: "installment",
  BANK_ACCOUNT: "bank_account",
  CREDIT_CARD: "credit_card",
  CASH: "cash",  // New
} as const;
```

#### 2. Form Validation
- Updated `createCashFlowSchema` to allow `target_type = "cash"` without requiring `target_account_id`
- Updated type conversion functions to map `PaymentMethodType.CASH` to `SourceType.CASH`

#### 3. UI Component (`frontend/src/components/cash-flows/AddCashFlowDialog.tsx`)
- Detects when "提領" category is selected
- Automatically sets `target_payment_method = "cash"`
- Hides transfer target fields when withdrawal is selected
- Properly submits data with `target_type = "cash"` and `target_id = null`

## Test Scripts

### Updated Files
- `backend/scripts/seed_test_categories.sql`: Added withdrawal category to test seed data

## Usage

### Creating a Cash Withdrawal Record

1. Open the cash flow page
2. Click "新增記錄" (Add Record)
3. Select type "轉出" (Transfer Out)
4. Select category "提領" (Withdrawal)
5. Fill in required fields:
   - Date
   - Amount
   - Description (e.g., "ATM 提領現金")
   - Bank Account (select the account to withdraw from)
6. Submit the form

### Expected Behavior

- Bank account balance is reduced by the withdrawal amount
- Cash flow record is created with:
  - `type = "transfer_out"`
  - `category_id = <withdrawal category ID>`
  - `source_type = "bank_account"`
  - `source_id = <bank account ID>`
  - `target_type = "cash"`
  - `target_id = null`

## Deployment

The migration files will automatically run when deploying to production using:

```bash
make migrate-up-env
```

No manual database operations are required.

## Testing

All tests pass successfully:
- Backend unit tests: ✅
- Backend service tests: ✅
- Backend repository tests: ✅
- All existing tests: ✅ (200+ tests)

## Future Enhancements

Potential improvements for future iterations:
- Add "Quick Withdrawal" button in the UI
- Add withdrawal statistics in reports
- Add filtering for withdrawal records
- Track ATM fees separately

