import { Instrument } from "@/types/instrument";
import { AssetType } from "@/types/transaction";

/**
 * 常用台股標的清單（台股市值前 30 大）
 */
const COMMON_TW_STOCKS: Instrument[] = [
  { symbol: "2330", name: "台積電", asset_type: AssetType.TW_STOCK },
  { symbol: "2317", name: "鴻海", asset_type: AssetType.TW_STOCK },
  { symbol: "2454", name: "聯發科", asset_type: AssetType.TW_STOCK },
  { symbol: "2308", name: "台達電", asset_type: AssetType.TW_STOCK },
  { symbol: "2882", name: "國泰金", asset_type: AssetType.TW_STOCK },
  { symbol: "2881", name: "富邦金", asset_type: AssetType.TW_STOCK },
  { symbol: "2886", name: "兆豐金", asset_type: AssetType.TW_STOCK },
  { symbol: "2891", name: "中信金", asset_type: AssetType.TW_STOCK },
  { symbol: "2412", name: "中華電", asset_type: AssetType.TW_STOCK },
  { symbol: "2303", name: "聯電", asset_type: AssetType.TW_STOCK },
  { symbol: "1301", name: "台塑", asset_type: AssetType.TW_STOCK },
  { symbol: "1303", name: "南亞", asset_type: AssetType.TW_STOCK },
  { symbol: "2002", name: "中鋼", asset_type: AssetType.TW_STOCK },
  { symbol: "2884", name: "玉山金", asset_type: AssetType.TW_STOCK },
  { symbol: "2892", name: "第一金", asset_type: AssetType.TW_STOCK },
  { symbol: "2912", name: "統一超", asset_type: AssetType.TW_STOCK },
  { symbol: "2357", name: "華碩", asset_type: AssetType.TW_STOCK },
  { symbol: "2382", name: "廣達", asset_type: AssetType.TW_STOCK },
  { symbol: "3008", name: "大立光", asset_type: AssetType.TW_STOCK },
  { symbol: "2327", name: "國巨", asset_type: AssetType.TW_STOCK },
  { symbol: "2395", name: "研華", asset_type: AssetType.TW_STOCK },
  { symbol: "2408", name: "南亞科", asset_type: AssetType.TW_STOCK },
  { symbol: "3711", name: "日月光投控", asset_type: AssetType.TW_STOCK },
  { symbol: "2603", name: "長榮", asset_type: AssetType.TW_STOCK },
  { symbol: "2609", name: "陽明", asset_type: AssetType.TW_STOCK },
  { symbol: "2615", name: "萬海", asset_type: AssetType.TW_STOCK },
  { symbol: "2207", name: "和泰車", asset_type: AssetType.TW_STOCK },
  { symbol: "2301", name: "光寶科", asset_type: AssetType.TW_STOCK },
  { symbol: "2379", name: "瑞昱", asset_type: AssetType.TW_STOCK },
  { symbol: "6505", name: "台塑化", asset_type: AssetType.TW_STOCK },
];

/**
 * 常用美股標的清單（S&P 100 精選）
 */
const COMMON_US_STOCKS: Instrument[] = [
  { symbol: "AAPL", name: "Apple Inc", asset_type: AssetType.US_STOCK },
  { symbol: "MSFT", name: "Microsoft Corporation", asset_type: AssetType.US_STOCK },
  { symbol: "GOOGL", name: "Alphabet Inc Class A", asset_type: AssetType.US_STOCK },
  { symbol: "AMZN", name: "Amazon.com Inc", asset_type: AssetType.US_STOCK },
  { symbol: "NVDA", name: "NVIDIA Corporation", asset_type: AssetType.US_STOCK },
  { symbol: "META", name: "Meta Platforms Inc", asset_type: AssetType.US_STOCK },
  { symbol: "TSLA", name: "Tesla Inc", asset_type: AssetType.US_STOCK },
  { symbol: "BRK.B", name: "Berkshire Hathaway Inc Class B", asset_type: AssetType.US_STOCK },
  { symbol: "V", name: "Visa Inc", asset_type: AssetType.US_STOCK },
  { symbol: "JNJ", name: "Johnson & Johnson", asset_type: AssetType.US_STOCK },
  { symbol: "WMT", name: "Walmart Inc", asset_type: AssetType.US_STOCK },
  { symbol: "JPM", name: "JPMorgan Chase & Co", asset_type: AssetType.US_STOCK },
  { symbol: "MA", name: "Mastercard Inc", asset_type: AssetType.US_STOCK },
  { symbol: "PG", name: "Procter & Gamble Co", asset_type: AssetType.US_STOCK },
  { symbol: "UNH", name: "UnitedHealth Group Inc", asset_type: AssetType.US_STOCK },
  { symbol: "HD", name: "Home Depot Inc", asset_type: AssetType.US_STOCK },
  { symbol: "DIS", name: "Walt Disney Co", asset_type: AssetType.US_STOCK },
  { symbol: "BAC", name: "Bank of America Corp", asset_type: AssetType.US_STOCK },
  { symbol: "ADBE", name: "Adobe Inc", asset_type: AssetType.US_STOCK },
  { symbol: "CRM", name: "Salesforce Inc", asset_type: AssetType.US_STOCK },
  { symbol: "NFLX", name: "Netflix Inc", asset_type: AssetType.US_STOCK },
  { symbol: "CSCO", name: "Cisco Systems Inc", asset_type: AssetType.US_STOCK },
  { symbol: "INTC", name: "Intel Corporation", asset_type: AssetType.US_STOCK },
  { symbol: "PEP", name: "PepsiCo Inc", asset_type: AssetType.US_STOCK },
  { symbol: "KO", name: "Coca-Cola Co", asset_type: AssetType.US_STOCK },
];

/**
 * 常用加密貨幣標的清單（市值前 20 大）
 */
const COMMON_CRYPTO: Instrument[] = [
  { symbol: "BTC", name: "Bitcoin", asset_type: AssetType.CRYPTO },
  { symbol: "ETH", name: "Ethereum", asset_type: AssetType.CRYPTO },
  { symbol: "USDT", name: "Tether", asset_type: AssetType.CRYPTO },
  { symbol: "BNB", name: "Binance Coin", asset_type: AssetType.CRYPTO },
  { symbol: "SOL", name: "Solana", asset_type: AssetType.CRYPTO },
  { symbol: "USDC", name: "USD Coin", asset_type: AssetType.CRYPTO },
  { symbol: "XRP", name: "Ripple", asset_type: AssetType.CRYPTO },
  { symbol: "ADA", name: "Cardano", asset_type: AssetType.CRYPTO },
  { symbol: "DOGE", name: "Dogecoin", asset_type: AssetType.CRYPTO },
  { symbol: "TRX", name: "TRON", asset_type: AssetType.CRYPTO },
  { symbol: "AVAX", name: "Avalanche", asset_type: AssetType.CRYPTO },
  { symbol: "DOT", name: "Polkadot", asset_type: AssetType.CRYPTO },
  { symbol: "MATIC", name: "Polygon", asset_type: AssetType.CRYPTO },
  { symbol: "LTC", name: "Litecoin", asset_type: AssetType.CRYPTO },
  { symbol: "LINK", name: "Chainlink", asset_type: AssetType.CRYPTO },
  { symbol: "UNI", name: "Uniswap", asset_type: AssetType.CRYPTO },
  { symbol: "ATOM", name: "Cosmos", asset_type: AssetType.CRYPTO },
  { symbol: "XLM", name: "Stellar", asset_type: AssetType.CRYPTO },
  { symbol: "ALGO", name: "Algorand", asset_type: AssetType.CRYPTO },
  { symbol: "VET", name: "VeChain", asset_type: AssetType.CRYPTO },
];

/**
 * 所有常用標的清單（合併）
 */
export const COMMON_INSTRUMENTS: Instrument[] = [
  ...COMMON_TW_STOCKS,
  ...COMMON_US_STOCKS,
  ...COMMON_CRYPTO,
];

/**
 * 根據資產類型取得常用標的清單
 * 
 * @param assetType 資產類型
 * @returns 該資產類型的常用標的清單
 * 
 * @example
 * ```typescript
 * const twStocks = getCommonInstrumentsByType(AssetType.TW_STOCK);
 * // 回傳台股常用標的清單
 * ```
 */
export function getCommonInstrumentsByType(assetType: AssetType): Instrument[] {
  return COMMON_INSTRUMENTS.filter((instrument) => instrument.asset_type === assetType);
}

/**
 * 搜尋常用標的（模糊搜尋 symbol 或 name）
 * 
 * @param query 搜尋關鍵字
 * @param assetType 資產類型（可選）
 * @returns 符合條件的標的清單
 * 
 * @example
 * ```typescript
 * const results = searchCommonInstruments("台積", AssetType.TW_STOCK);
 * // 回傳 [{ symbol: "2330", name: "台積電", asset_type: "tw-stock" }]
 * ```
 */
export function searchCommonInstruments(
  query: string,
  assetType?: AssetType
): Instrument[] {
  const normalizedQuery = query.toLowerCase().trim();
  
  let instruments = COMMON_INSTRUMENTS;
  
  // 如果指定資產類型，先篩選
  if (assetType) {
    instruments = instruments.filter((i) => i.asset_type === assetType);
  }
  
  // 模糊搜尋 symbol 或 name
  return instruments.filter(
    (instrument) =>
      instrument.symbol.toLowerCase().includes(normalizedQuery) ||
      instrument.name.toLowerCase().includes(normalizedQuery)
  );
}

