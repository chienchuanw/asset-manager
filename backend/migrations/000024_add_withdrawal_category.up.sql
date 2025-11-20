-- Add withdrawal category for cash withdrawals from bank accounts
-- This category is used when users withdraw cash from their bank accounts

-- Add withdrawal category for transfer_out type
INSERT INTO cash_flow_categories (name, type, is_system) 
VALUES ('提領', 'transfer_out', true)
ON CONFLICT (name, type) DO NOTHING;

-- Add comment
COMMENT ON TABLE cash_flow_categories IS 'Cash flow categories table. Includes system categories: income, expense, transfer_in, transfer_out, and withdrawal';

