-- Update cash flow source type constraint to include 'cash' type
-- This is used for cash withdrawal target type

-- Remove existing constraint
ALTER TABLE cash_flows DROP CONSTRAINT IF EXISTS cash_flows_source_type_check;

-- Add updated constraint including 'cash' type
ALTER TABLE cash_flows 
ADD CONSTRAINT cash_flows_source_type_check 
CHECK (source_type IN ('manual', 'subscription', 'installment', 'bank_account', 'credit_card', 'cash'));

-- Update comment
COMMENT ON COLUMN cash_flows.source_type IS 'Source type (manual: manual/cash transaction, subscription: auto-generated from subscription, installment: auto-generated from installment, bank_account: bank account transaction, credit_card: credit card transaction, cash: cash)';
COMMENT ON COLUMN cash_flows.source_id IS 'Source ID (references subscription, installment, bank account, or credit card ID)';

-- Update target_type comment to include cash
COMMENT ON COLUMN cash_flows.target_type IS 'Transfer target type (used for transfer_in/transfer_out, e.g., credit_card, bank_account, cash)';

