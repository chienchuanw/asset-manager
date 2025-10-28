import { AssetType } from "./transaction";

/**
 * 標的（股票/加密貨幣）資料結構
 */
export interface Instrument {
  /** 標的代碼（例如：2330, AAPL, BTC） */
  symbol: string;
  /** 標的名稱（例如：台積電, Apple Inc, Bitcoin） */
  name: string;
  /** 資產類型 */
  asset_type: AssetType;
}

/**
 * 標的搜尋結果（包含額外資訊）
 */
export interface InstrumentSearchResult extends Instrument {
  /** 是否來自歷史交易記錄 */
  from_history?: boolean;
  /** 使用次數（用於排序） */
  usage_count?: number;
}

