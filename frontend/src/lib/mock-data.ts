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
  type: 'buy' | 'sell';
  quantity: number;
  price: number;
  amount: number;
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

// 近期交易 Mock Data
export const mockTransactions: Transaction[] = [
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

