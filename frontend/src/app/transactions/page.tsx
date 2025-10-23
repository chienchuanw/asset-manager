/**
 * 交易記錄頁面
 * 顯示所有交易記錄,支援篩選、搜尋、排序功能
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
import {
  mockAllTransactions,
  assetTypeNames,
  transactionTypeNames,
  transactionTypeColors,
  assetTypeColors,
  type AssetType,
  type Transaction,
} from '@/lib/mock-data';
import { Search, Download, Plus, ArrowUpDown } from 'lucide-react';

export default function TransactionsPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [filterType, setFilterType] = useState<Transaction['type'] | 'all'>('all');
  const [filterAssetType, setFilterAssetType] = useState<AssetType | 'all'>('all');
  const [sortBy, setSortBy] = useState<'date' | 'amount'>('date');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  // 篩選和排序邏輯
  const filteredAndSortedTransactions = mockAllTransactions
    .filter((transaction) => {
      // 交易類型篩選
      if (filterType !== 'all' && transaction.type !== filterType) {
        return false;
      }
      // 資產類別篩選
      if (filterAssetType !== 'all' && transaction.assetType !== filterAssetType) {
        return false;
      }
      // 搜尋篩選
      if (searchQuery) {
        const query = searchQuery.toLowerCase();
        return (
          transaction.symbol.toLowerCase().includes(query) ||
          transaction.name.toLowerCase().includes(query)
        );
      }
      return true;
    })
    .sort((a, b) => {
      let compareValue = 0;
      switch (sortBy) {
        case 'date':
          compareValue = new Date(a.date).getTime() - new Date(b.date).getTime();
          break;
        case 'amount':
          compareValue = a.amount - b.amount;
          break;
      }
      return sortOrder === 'asc' ? compareValue : -compareValue;
    });

  // 計算統計資料
  const totalTransactions = filteredAndSortedTransactions.length;
  const totalBuyAmount = filteredAndSortedTransactions
    .filter((t) => t.type === 'buy')
    .reduce((sum, t) => sum + t.amount, 0);
  const totalSellAmount = filteredAndSortedTransactions
    .filter((t) => t.type === 'sell')
    .reduce((sum, t) => sum + t.amount, 0);
  const netFlow = totalBuyAmount - totalSellAmount;

  // 排序切換
  const handleSort = (field: 'date' | 'amount') => {
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
        <div className="flex flex-col gap-6">
          {/* 統計摘要卡片 */}
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>總交易次數</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold tabular-nums">{totalTransactions}</div>
                <p className="text-xs text-muted-foreground mt-1">筆交易記錄</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardDescription>總買入金額</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold tabular-nums text-blue-600">
                  TWD {totalBuyAmount.toLocaleString()}
                </div>
                <p className="text-xs text-muted-foreground mt-1">累計買入</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardDescription>總賣出金額</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold tabular-nums text-red-600">
                  TWD {totalSellAmount.toLocaleString()}
                </div>
                <p className="text-xs text-muted-foreground mt-1">累計賣出</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardDescription>淨流入/流出</CardDescription>
              </CardHeader>
              <CardContent>
                <div
                  className={`text-2xl font-bold tabular-nums ${
                    netFlow >= 0 ? 'text-green-600' : 'text-red-600'
                  }`}
                >
                  TWD {netFlow.toLocaleString()}
                </div>
                <p className="text-xs text-muted-foreground mt-1">
                  {netFlow >= 0 ? '淨流入' : '淨流出'}
                </p>
              </CardContent>
            </Card>
          </div>

          {/* 交易記錄表格 */}
          <Card>
            <CardHeader>
              <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <CardTitle>交易記錄</CardTitle>
                  <CardDescription>查看所有交易記錄</CardDescription>
                </div>
                <div className="flex gap-2">
                  <Button variant="outline" size="sm">
                    <Plus className="h-4 w-4 mr-2" />
                    新增交易
                  </Button>
                  <Button variant="outline" size="sm">
                    <Download className="h-4 w-4 mr-2" />
                    匯出
                  </Button>
                </div>
              </div>

              {/* 篩選工具列 */}
              <div className="flex flex-col gap-3 sm:flex-row sm:items-center mt-4">
                {/* 搜尋框 */}
                <div className="relative flex-1">
                  <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                  <Input
                    placeholder="搜尋代碼或名稱..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="pl-9"
                  />
                </div>

                {/* 交易類型篩選 */}
                <Select value={filterType} onValueChange={(value) => setFilterType(value as any)}>
                  <SelectTrigger className="w-full sm:w-[150px]">
                    <SelectValue placeholder="交易類型" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">全部類型</SelectItem>
                    <SelectItem value="buy">買入</SelectItem>
                    <SelectItem value="sell">賣出</SelectItem>
                    <SelectItem value="dividend">股利</SelectItem>
                    <SelectItem value="fee">手續費</SelectItem>
                  </SelectContent>
                </Select>

                {/* 資產類別篩選 */}
                <Select
                  value={filterAssetType}
                  onValueChange={(value) => setFilterAssetType(value as any)}
                >
                  <SelectTrigger className="w-full sm:w-[150px]">
                    <SelectValue placeholder="資產類別" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">全部類別</SelectItem>
                    <SelectItem value="tw-stock">台股</SelectItem>
                    <SelectItem value="us-stock">美股</SelectItem>
                    <SelectItem value="crypto">加密貨幣</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </CardHeader>

            <CardContent>
              <div className="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleSort('date')}
                          className="h-8 px-2"
                        >
                          日期
                          <ArrowUpDown className="ml-2 h-4 w-4" />
                        </Button>
                      </TableHead>
                      <TableHead>交易類型</TableHead>
                      <TableHead>資產類別</TableHead>
                      <TableHead>代碼/名稱</TableHead>
                      <TableHead className="text-right">數量</TableHead>
                      <TableHead className="text-right">單價</TableHead>
                      <TableHead>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleSort('amount')}
                          className="h-8 px-2 ml-auto"
                        >
                          總金額
                          <ArrowUpDown className="ml-2 h-4 w-4" />
                        </Button>
                      </TableHead>
                      <TableHead className="text-right hidden md:table-cell">手續費</TableHead>
                      <TableHead className="hidden lg:table-cell">備註</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {filteredAndSortedTransactions.map((transaction) => (
                      <TableRow key={transaction.id}>
                        <TableCell className="font-medium">{transaction.date}</TableCell>
                        <TableCell>
                          <Badge variant="outline" className={transactionTypeColors[transaction.type]}>
                            {transactionTypeNames[transaction.type]}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <Badge variant="outline" className={assetTypeColors[transaction.assetType]}>
                            {assetTypeNames[transaction.assetType]}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <div>
                            <div className="font-medium">{transaction.symbol}</div>
                            <div className="text-sm text-muted-foreground">{transaction.name}</div>
                          </div>
                        </TableCell>
                        <TableCell className="text-right tabular-nums">
                          {transaction.quantity.toLocaleString()}
                        </TableCell>
                        <TableCell className="text-right tabular-nums">
                          {transaction.price.toLocaleString()}
                        </TableCell>
                        <TableCell className="text-right font-medium tabular-nums">
                          {transaction.amount.toLocaleString()}
                        </TableCell>
                        <TableCell className="text-right tabular-nums hidden md:table-cell">
                          {transaction.fee ? transaction.fee.toLocaleString() : '-'}
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground hidden lg:table-cell">
                          {transaction.note || '-'}
                        </TableCell>
                      </TableRow>
                    ))}
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

