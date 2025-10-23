/**
 * 持倉明細頁面
 * 顯示所有資產的詳細持倉資訊,支援篩選、排序、搜尋功能
 */

'use client';

import { useState } from 'react';
import { AppLayout } from '@/components/layout/AppLayout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Button } from '@/components/ui/button';
import { mockHoldings, assetTypeNames, type AssetType } from '@/lib/mock-data';
import { Search, ArrowUpDown, Download } from 'lucide-react';

export default function HoldingsPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [filterType, setFilterType] = useState<AssetType | 'all'>('all');
  const [sortBy, setSortBy] = useState<'value' | 'profit' | 'quantity'>('value');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  // 篩選與排序邏輯
  const filteredAndSortedHoldings = mockHoldings
    .filter((holding) => {
      // 資產類別篩選
      if (filterType !== 'all' && holding.assetType !== filterType) {
        return false;
      }
      // 搜尋篩選
      if (searchQuery) {
        const query = searchQuery.toLowerCase();
        return (
          holding.symbol.toLowerCase().includes(query) ||
          holding.name.toLowerCase().includes(query)
        );
      }
      return true;
    })
    .sort((a, b) => {
      let compareValue = 0;
      switch (sortBy) {
        case 'value':
          compareValue = a.marketValue - b.marketValue;
          break;
        case 'profit':
          compareValue = a.profitLoss - b.profitLoss;
          break;
        case 'quantity':
          compareValue = a.quantity - b.quantity;
          break;
      }
      return sortOrder === 'asc' ? compareValue : -compareValue;
    });

  // 計算統計資料
  const totalMarketValue = filteredAndSortedHoldings.reduce(
    (sum, h) => sum + h.marketValue,
    0
  );
  const totalCost = filteredAndSortedHoldings.reduce((sum, h) => sum + h.cost, 0);
  const totalProfitLoss = filteredAndSortedHoldings.reduce(
    (sum, h) => sum + h.profitLoss,
    0
  );
  const totalProfitLossPercent =
    totalCost > 0 ? (totalProfitLoss / totalCost) * 100 : 0;

  // 切換排序
  const toggleSort = (field: 'value' | 'profit' | 'quantity') => {
    if (sortBy === field) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(field);
      setSortOrder('desc');
    }
  };

  return (
    <AppLayout>
      {/* Main Content */}
      <main className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="container flex flex-col gap-6">
          {/* 統計摘要卡片 */}
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>總市值</CardDescription>
                <CardTitle className="text-2xl tabular-nums">
                  NT$ {totalMarketValue.toLocaleString()}
                </CardTitle>
              </CardHeader>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>總成本</CardDescription>
                <CardTitle className="text-2xl tabular-nums">
                  NT$ {totalCost.toLocaleString()}
                </CardTitle>
              </CardHeader>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>未實現損益</CardDescription>
                <CardTitle
                  className={`text-2xl tabular-nums ${
                    totalProfitLoss >= 0 ? 'text-green-600' : 'text-red-600'
                  }`}
                >
                  {totalProfitLoss >= 0 ? '+' : ''}NT${' '}
                  {totalProfitLoss.toLocaleString()}
                </CardTitle>
              </CardHeader>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>報酬率</CardDescription>
                <CardTitle
                  className={`text-2xl tabular-nums ${
                    totalProfitLossPercent >= 0 ? 'text-green-600' : 'text-red-600'
                  }`}
                >
                  {totalProfitLossPercent >= 0 ? '+' : ''}
                  {totalProfitLossPercent.toFixed(2)}%
                </CardTitle>
              </CardHeader>
            </Card>
          </div>

          {/* 篩選與搜尋工具列 */}
          <Card>
            <CardHeader>
              <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <CardTitle>持倉列表</CardTitle>
                  <CardDescription>
                    共 {filteredAndSortedHoldings.length} 筆持倉
                  </CardDescription>
                </div>
                <div className="flex flex-col gap-2 sm:flex-row">
                  {/* 搜尋框 */}
                  <div className="relative">
                    <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input
                      placeholder="搜尋代碼或名稱..."
                      className="pl-8 sm:w-[200px]"
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                    />
                  </div>

                  {/* 資產類別篩選 */}
                  <Select
                    value={filterType}
                    onValueChange={(value) => setFilterType(value as AssetType | 'all')}
                  >
                    <SelectTrigger className="sm:w-[150px]">
                      <SelectValue placeholder="資產類別" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">全部類別</SelectItem>
                      <SelectItem value="cash">現金</SelectItem>
                      <SelectItem value="tw-stock">台股</SelectItem>
                      <SelectItem value="us-stock">美股</SelectItem>
                      <SelectItem value="crypto">加密貨幣</SelectItem>
                    </SelectContent>
                  </Select>

                  {/* 匯出按鈕 */}
                  <Button variant="outline" size="sm">
                    <Download className="mr-2 h-4 w-4" />
                    匯出
                  </Button>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              {/* 持倉表格 */}
              <div className="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>代碼/名稱</TableHead>
                      <TableHead>類別</TableHead>
                      <TableHead
                        className="cursor-pointer"
                        onClick={() => toggleSort('quantity')}
                      >
                        <div className="flex items-center gap-1">
                          持有數量
                          <ArrowUpDown className="h-3 w-3" />
                        </div>
                      </TableHead>
                      <TableHead className="text-right">成本價</TableHead>
                      <TableHead className="text-right">現價</TableHead>
                      <TableHead
                        className="cursor-pointer text-right"
                        onClick={() => toggleSort('value')}
                      >
                        <div className="flex items-center justify-end gap-1">
                          市值
                          <ArrowUpDown className="h-3 w-3" />
                        </div>
                      </TableHead>
                      <TableHead
                        className="cursor-pointer text-right"
                        onClick={() => toggleSort('profit')}
                      >
                        <div className="flex items-center justify-end gap-1">
                          損益
                          <ArrowUpDown className="h-3 w-3" />
                        </div>
                      </TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {filteredAndSortedHoldings.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={7} className="h-24 text-center">
                          沒有符合條件的持倉
                        </TableCell>
                      </TableRow>
                    ) : (
                      filteredAndSortedHoldings.map((holding) => (
                        <TableRow key={holding.id}>
                          <TableCell>
                            <div className="flex flex-col">
                              <span className="font-medium">{holding.symbol}</span>
                              <span className="text-sm text-muted-foreground">
                                {holding.name}
                              </span>
                            </div>
                          </TableCell>
                          <TableCell>
                            <Badge variant="outline">
                              {assetTypeNames[holding.assetType]}
                            </Badge>
                          </TableCell>
                          <TableCell className="tabular-nums">
                            {holding.quantity.toLocaleString()}
                          </TableCell>
                          <TableCell className="text-right tabular-nums">
                            {holding.avgCost.toLocaleString()}
                          </TableCell>
                          <TableCell className="text-right tabular-nums">
                            {holding.currentPrice.toLocaleString()}
                          </TableCell>
                          <TableCell className="text-right tabular-nums">
                            {holding.marketValue.toLocaleString()}
                          </TableCell>
                          <TableCell className="text-right">
                            <div className="flex flex-col items-end">
                              <span
                                className={`font-medium tabular-nums ${
                                  holding.profitLoss >= 0
                                    ? 'text-green-600'
                                    : 'text-red-600'
                                }`}
                              >
                                {holding.profitLoss >= 0 ? '+' : ''}
                                {holding.profitLoss.toLocaleString()}
                              </span>
                              <span
                                className={`text-sm tabular-nums ${
                                  holding.profitLossPercent >= 0
                                    ? 'text-green-600'
                                    : 'text-red-600'
                                }`}
                              >
                                {holding.profitLossPercent >= 0 ? '+' : ''}
                                {holding.profitLossPercent.toFixed(2)}%
                              </span>
                            </div>
                          </TableCell>
                        </TableRow>
                      ))
                    )}
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </div>
      </main>
    </AppLayout>
  );
}

