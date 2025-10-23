/**
 * Mock Data for Asset Manager Dashboard
 * 用於 UI 展示的模擬資料
 */

// 資產類型定義
export type AssetType = 'cash' | 'tw-stock' | 'us-stock' | 'crypto';

// 統計卡片資料
export interface StatCard {
  title: string;
  value: string;
  change: number; // 百分比變化
  description?: string; // 額外說明文字
}

// 持倉資料
export interface Holding {
  id: string;
  assetType: AssetType;
  symbol: string;
  name: string;
  quantity: number;
  avgCost: number; // 平均成本價
  cost: number; // 總成本 (avgCost * quantity)
  currentPrice: number;
  marketValue: number;
  profitLoss: number;
  profitLossPercent: number;
}

// 交易記錄
export interface Transaction {
  id: string;
  date: string;
  assetType: AssetType;
  symbol: string;
  name: string;
  type: 'buy' | 'sell' | 'dividend' | 'fee';
  quantity: number;
  price: number;
  amount: number;
  fee?: number; // 手續費
  note?: string; // 備註
}

// 資產配置
export interface AssetAllocation {
  assetType: AssetType;
  name: string;
  value: number;
  percentage: number;
  color: string;
}

// 圖表資料點
export interface ChartDataPoint {
  date: string;
  total: number;
  cash: number;
  twStock: number;
  usStock: number;
  crypto: number;
}

// 統計卡片 Mock Data
export const mockStatCards: StatCard[] = [
  {
    title: '總資產價值',
    value: 'NT$ 1,234,567',
    change: 5.2,
    description: '本月資產持續成長',
  },
  {
    title: '今日損益',
    value: 'NT$ 12,345',
    change: 2.3,
    description: '今日表現優於大盤',
  },
  {
    title: '持倉數量',
    value: '24',
    change: 8.0,
    description: '持倉組合多元化',
  },
  {
    title: '總報酬率',
    value: '18.5%',
    change: 3.2,
    description: '年化報酬率穩定成長',
  },
];

// 持倉明細 Mock Data
export const mockHoldings: Holding[] = [
  {
    id: '1',
    assetType: 'tw-stock',
    symbol: '2330',
    name: '台積電',
    quantity: 100,
    avgCost: 580,
    cost: 58000,
    currentPrice: 620,
    marketValue: 62000,
    profitLoss: 4000,
    profitLossPercent: 6.9,
  },
  {
    id: '2',
    assetType: 'us-stock',
    symbol: 'AAPL',
    name: 'Apple Inc.',
    quantity: 50,
    avgCost: 150,
    cost: 7500,
    currentPrice: 175,
    marketValue: 8750,
    profitLoss: 1250,
    profitLossPercent: 16.7,
  },
  {
    id: '3',
    assetType: 'crypto',
    symbol: 'BTC',
    name: 'Bitcoin',
    quantity: 0.5,
    avgCost: 900000,
    cost: 450000,
    currentPrice: 1200000,
    marketValue: 600000,
    profitLoss: 150000,
    profitLossPercent: 33.3,
  },
  {
    id: '4',
    assetType: 'tw-stock',
    symbol: '2317',
    name: '鴻海',
    quantity: 200,
    avgCost: 105,
    cost: 21000,
    currentPrice: 98,
    marketValue: 19600,
    profitLoss: -1400,
    profitLossPercent: -6.7,
  },
  {
    id: '5',
    assetType: 'us-stock',
    symbol: 'TSLA',
    name: 'Tesla Inc.',
    quantity: 30,
    avgCost: 250,
    cost: 7500,
    currentPrice: 280,
    marketValue: 8400,
    profitLoss: 900,
    profitLossPercent: 12.0,
  },
];

// 近期交易 Mock Data (用於 Dashboard)
export const mockRecentTransactions: Transaction[] = [
  {
    id: '1',
    date: '2025-10-22',
    assetType: 'tw-stock',
    symbol: '2330',
    name: '台積電',
    type: 'buy',
    quantity: 10,
    price: 620,
    amount: 6200,
    fee: 28,
  },
  {
    id: '2',
    date: '2025-10-21',
    assetType: 'crypto',
    symbol: 'ETH',
    name: 'Ethereum',
    type: 'buy',
    quantity: 2,
    price: 50000,
    amount: 100000,
    fee: 100,
  },
  {
    id: '3',
    date: '2025-10-20',
    assetType: 'us-stock',
    symbol: 'AAPL',
    name: 'Apple Inc.',
    type: 'sell',
    quantity: 5,
    price: 175,
    amount: 875,
    fee: 4,
  },
  {
    id: '4',
    date: '2025-10-19',
    assetType: 'tw-stock',
    symbol: '2317',
    name: '鴻海',
    type: 'buy',
    quantity: 50,
    price: 98,
    amount: 4900,
    fee: 22,
  },
  {
    id: '5',
    date: '2025-10-18',
    assetType: 'us-stock',
    symbol: 'TSLA',
    name: 'Tesla Inc.',
    type: 'buy',
    quantity: 10,
    price: 280,
    amount: 2800,
    fee: 14,
  },
];

// 完整交易記錄 Mock Data (用於交易記錄頁面)
export const mockAllTransactions: Transaction[] = [
  {
    id: '1',
    date: '2025-10-22',
    assetType: 'tw-stock',
    symbol: '2330',
    name: '台積電',
    type: 'buy',
    quantity: 10,
    price: 620,
    amount: 6200,
    fee: 28,
    note: '定期定額買入',
  },
  {
    id: '2',
    date: '2025-10-21',
    assetType: 'crypto',
    symbol: 'ETH',
    name: 'Ethereum',
    type: 'buy',
    quantity: 2,
    price: 50000,
    amount: 100000,
    fee: 100,
  },
  {
    id: '3',
    date: '2025-10-20',
    assetType: 'us-stock',
    symbol: 'AAPL',
    name: 'Apple Inc.',
    type: 'sell',
    quantity: 5,
    price: 175,
    amount: 875,
    fee: 4,
    note: '部分獲利了結',
  },
  {
    id: '4',
    date: '2025-10-19',
    assetType: 'tw-stock',
    symbol: '2317',
    name: '鴻海',
    type: 'buy',
    quantity: 50,
    price: 98,
    amount: 4900,
    fee: 22,
  },
  {
    id: '5',
    date: '2025-10-18',
    assetType: 'us-stock',
    symbol: 'TSLA',
    name: 'Tesla Inc.',
    type: 'buy',
    quantity: 10,
    price: 280,
    amount: 2800,
    fee: 14,
  },
  {
    id: '6',
    date: '2025-10-15',
    assetType: 'tw-stock',
    symbol: '2330',
    name: '台積電',
    type: 'dividend',
    quantity: 100,
    price: 2.75,
    amount: 275,
    note: '現金股利',
  },
  {
    id: '7',
    date: '2025-10-12',
    assetType: 'crypto',
    symbol: 'BTC',
    name: 'Bitcoin',
    type: 'buy',
    quantity: 0.5,
    price: 900000,
    amount: 450000,
    fee: 450,
  },
  {
    id: '8',
    date: '2025-10-10',
    assetType: 'us-stock',
    symbol: 'AAPL',
    name: 'Apple Inc.',
    type: 'buy',
    quantity: 20,
    price: 150,
    amount: 3000,
    fee: 15,
  },
  {
    id: '9',
    date: '2025-10-08',
    assetType: 'tw-stock',
    symbol: '2317',
    name: '鴻海',
    type: 'buy',
    quantity: 150,
    price: 105,
    amount: 15750,
    fee: 71,
  },
  {
    id: '10',
    date: '2025-10-05',
    assetType: 'us-stock',
    symbol: 'TSLA',
    name: 'Tesla Inc.',
    type: 'buy',
    quantity: 20,
    price: 250,
    amount: 5000,
    fee: 25,
  },
  {
    id: '11',
    date: '2025-10-01',
    assetType: 'tw-stock',
    symbol: '2330',
    name: '台積電',
    type: 'buy',
    quantity: 90,
    price: 580,
    amount: 52200,
    fee: 235,
    note: '大量買入',
  },
  {
    id: '12',
    date: '2025-09-28',
    assetType: 'us-stock',
    symbol: 'AAPL',
    name: 'Apple Inc.',
    type: 'buy',
    quantity: 30,
    price: 148,
    amount: 4440,
    fee: 22,
  },
];

// 資產配置 Mock Data
export const mockAssetAllocation: AssetAllocation[] = [
  {
    assetType: 'cash',
    name: '現金',
    value: 300000,
    percentage: 24.3,
    color: '#10b981', // green
  },
  {
    assetType: 'tw-stock',
    name: '台股',
    value: 400000,
    percentage: 32.4,
    color: '#3b82f6', // blue
  },
  {
    assetType: 'us-stock',
    name: '美股',
    value: 350000,
    percentage: 28.3,
    color: '#8b5cf6', // purple
  },
  {
    assetType: 'crypto',
    name: '加密貨幣',
    value: 184000,
    percentage: 15.0,
    color: '#f59e0b', // amber
  },
];

// 資產趨勢圖表 Mock Data (最近 30 天)
export const mockChartData: ChartDataPoint[] = [
  { date: '09/24', total: 1100000, cash: 300000, twStock: 380000, usStock: 320000, crypto: 100000 },
  { date: '09/27', total: 1120000, cash: 300000, twStock: 385000, usStock: 325000, crypto: 110000 },
  { date: '09/30', total: 1105000, cash: 300000, twStock: 375000, usStock: 330000, crypto: 100000 },
  { date: '10/03', total: 1150000, cash: 300000, twStock: 390000, usStock: 340000, crypto: 120000 },
  { date: '10/06', total: 1140000, cash: 300000, twStock: 385000, usStock: 335000, crypto: 120000 },
  { date: '10/09', total: 1180000, cash: 300000, twStock: 395000, usStock: 345000, crypto: 140000 },
  { date: '10/12', total: 1170000, cash: 300000, twStock: 390000, usStock: 340000, crypto: 140000 },
  { date: '10/15', total: 1200000, cash: 300000, twStock: 400000, usStock: 350000, crypto: 150000 },
  { date: '10/18', total: 1190000, cash: 300000, twStock: 395000, usStock: 345000, crypto: 150000 },
  { date: '10/21', total: 1220000, cash: 300000, twStock: 405000, usStock: 355000, crypto: 160000 },
  { date: '10/23', total: 1234567, cash: 300000, twStock: 400000, usStock: 350000, crypto: 184567 },
];

// 資產類型顯示名稱對應
export const assetTypeNames: Record<AssetType, string> = {
  cash: '現金',
  'tw-stock': '台股',
  'us-stock': '美股',
  crypto: '加密貨幣',
};

// 資產類型顏色對應
export const assetTypeColors: Record<AssetType, string> = {
  cash: 'bg-green-100 text-green-800',
  'tw-stock': 'bg-blue-100 text-blue-800',
  'us-stock': 'bg-purple-100 text-purple-800',
  crypto: 'bg-amber-100 text-amber-800',
};

// 交易類型名稱對應
export const transactionTypeNames: Record<Transaction['type'], string> = {
  buy: '買入',
  sell: '賣出',
  dividend: '股利',
  fee: '手續費',
};

// 交易類型顏色對應
export const transactionTypeColors: Record<Transaction['type'], string> = {
  buy: 'bg-blue-100 text-blue-800',
  sell: 'bg-red-100 text-red-800',
  dividend: 'bg-green-100 text-green-800',
  fee: 'bg-gray-100 text-gray-800',
};

// 績效分析資料
export interface PerformanceData {
  assetType: AssetType;
  name: string;
  returnRate: number; // 報酬率 (%)
  profit: number; // 損益金額
}

// 各資產報酬率 Mock Data
export const mockPerformanceData: PerformanceData[] = [
  {
    assetType: 'tw-stock',
    name: '台股',
    returnRate: 8.5,
    profit: 5600,
  },
  {
    assetType: 'us-stock',
    name: '美股',
    returnRate: 14.2,
    profit: 2150,
  },
  {
    assetType: 'crypto',
    name: '加密貨幣',
    returnRate: 33.3,
    profit: 150000,
  },
];

// Top 資產資料
export interface TopAsset {
  symbol: string;
  name: string;
  assetType: AssetType;
  profit: number;
  profitPercent: number;
}

// Top 5 獲利資產
export const mockTopProfitAssets: TopAsset[] = [
  {
    symbol: 'BTC',
    name: 'Bitcoin',
    assetType: 'crypto',
    profit: 150000,
    profitPercent: 33.3,
  },
  {
    symbol: '2330',
    name: '台積電',
    assetType: 'tw-stock',
    profit: 4000,
    profitPercent: 6.9,
  },
  {
    symbol: 'AAPL',
    name: 'Apple Inc.',
    assetType: 'us-stock',
    profit: 1250,
    profitPercent: 16.7,
  },
  {
    symbol: 'TSLA',
    name: 'Tesla Inc.',
    assetType: 'us-stock',
    profit: 900,
    profitPercent: 12.0,
  },
];

// Top 5 虧損資產
export const mockTopLossAssets: TopAsset[] = [
  {
    symbol: '2317',
    name: '鴻海',
    assetType: 'tw-stock',
    profit: -1400,
    profitPercent: -6.7,
  },
];
