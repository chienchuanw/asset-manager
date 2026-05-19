// FIRE (Financial Independence / Retire Early) 型別

/**
 * 投影圖表單一資料點
 */
export interface FireProjectionPoint {
  year: number;
  projected_net_worth: number;
  fire_target: number;
}

/**
 * FIRE 投影結果（後端回傳已套用的輸入值供前端回填）
 */
export interface FireProjectionResult {
  current_net_worth: number;
  annual_expenses: number;
  annual_savings: number;
  expected_return: number;
  withdrawal_rate: number;
  fire_number: number;
  progress_pct: number;
  years_to_fi: number | null;
  on_track: boolean;
  projection: FireProjectionPoint[];
}

/**
 * 可選的覆寫參數；省略的欄位由後端自動推導。
 */
export interface FireProjectionParams {
  netWorth?: number;
  annualExpenses?: number;
  annualSavings?: number;
  expectedReturn?: number;
  withdrawalRate?: number;
}
