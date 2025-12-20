-- 新增 sort_order 欄位到 cash_flow_categories 資料表
-- 用於支援分類的拖拉排序功能

-- 新增 sort_order 欄位，預設值為 0
ALTER TABLE cash_flow_categories
ADD COLUMN sort_order INTEGER NOT NULL DEFAULT 0;

-- 為現有分類設定初始排序值
-- 根據 type 分組，系統分類排在前面，按 created_at 排序
WITH ranked_categories AS (
    SELECT 
        id,
        type,
        ROW_NUMBER() OVER (
            PARTITION BY type 
            ORDER BY is_system DESC, created_at ASC
        ) - 1 AS new_order
    FROM cash_flow_categories
)
UPDATE cash_flow_categories c
SET sort_order = rc.new_order
FROM ranked_categories rc
WHERE c.id = rc.id;

-- 建立索引以提升排序查詢效能
CREATE INDEX idx_cash_flow_categories_sort_order ON cash_flow_categories(type, sort_order);

-- 新增註解說明
COMMENT ON COLUMN cash_flow_categories.sort_order IS '排序順序（0 開始，數字越小越前面）';

