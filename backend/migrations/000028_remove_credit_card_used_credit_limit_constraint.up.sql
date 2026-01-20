-- 移除信用卡已使用額度不能超過信用額度的約束
-- 此變更允許使用者手動調整已使用額度，用於對帳或修正系統計算錯誤
ALTER TABLE credit_cards DROP CONSTRAINT IF EXISTS check_used_credit_limit;
