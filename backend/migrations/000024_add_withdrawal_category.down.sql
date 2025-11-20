-- Rollback: Remove withdrawal category

-- Remove withdrawal category for transfer_out type
DELETE FROM cash_flow_categories 
WHERE name = '提領' AND type = 'transfer_out';

