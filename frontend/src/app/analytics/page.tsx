/**
 * 分析報表頁面
 * 顯示資產配置、績效分析、損益分析等報表
 */

'use client';

import { useState } from 'react';
import { AppLayout } from '@/components/layout/AppLayout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  mockAssetAllocation,
  mockPerformanceData,
  mockTopProfitAssets,
  mockTopLossAssets,
  mockChartData,
  assetTypeNames,
  assetTypeColors,
} from '@/lib/mock-data';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell } from 'recharts';
import { TrendingUp, TrendingDown } from 'lucide-react';

export default function AnalyticsPage() {
  const [timeRange, setTimeRange] = useState<'week' | 'month' | 'quarter' | 'year' | 'all'>('month');

  // 計算統計資料
  const totalValue = mockAssetAllocation.reduce((sum, item) => sum + item.value, 0);
  const totalProfit = mockPerformanceData.reduce((sum, item) => sum + item.profit, 0);
  const avgReturnRate =
    mockPerformanceData.reduce((sum, item) => sum + item.returnRate, 0) / mockPerformanceData.length;

  // 未實現損益 (從 holdings 計算)
  const unrealizedProfit = 157750; // Mock data
  const realizedProfit = 0; // Mock data (暫時為 0)

  return (
    <AppLayout>
      {/* Main Content */}
      <main className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="flex flex-col gap-6">
          {/* 時間範圍選擇 */}
          <Tabs value={timeRange} onValueChange={(value) => setTimeRange(value as any)}>
            <TabsList>
              <TabsTrigger value="week">本週</TabsTrigger>
              <TabsTrigger value="month">本月</TabsTrigger>
              <TabsTrigger value="quarter">本季</TabsTrigger>
              <TabsTrigger value="year">本年</TabsTrigger>
              <TabsTrigger value="all">全部</TabsTrigger>
            </TabsList>
          </Tabs>

          {/* 績效摘要卡片 */}
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>總資產價值</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold tabular-nums">
                  TWD {totalValue.toLocaleString()}
                </div>
                <p className="text-xs text-muted-foreground mt-1">當前總市值</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardDescription>總報酬率</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold tabular-nums text-green-600">
                  +{avgReturnRate.toFixed(2)}%
                </div>
                <p className="text-xs text-muted-foreground mt-1">
                  <Badge variant="outline" className="bg-green-100 text-green-800">
                    <TrendingUp className="h-3 w-3 mr-1" />
                    表現良好
                  </Badge>
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardDescription>未實現損益</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold tabular-nums text-green-600">
                  TWD {unrealizedProfit.toLocaleString()}
                </div>
                <p className="text-xs text-muted-foreground mt-1">帳面損益</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardDescription>已實現損益</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold tabular-nums text-gray-600">
                  TWD {realizedProfit.toLocaleString()}
                </div>
                <p className="text-xs text-muted-foreground mt-1">實際獲利</p>
              </CardContent>
            </Card>
          </div>

          {/* 圖表區域 */}
          <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
            {/* 各資產報酬率長條圖 */}
            <Card>
              <CardHeader>
                <CardTitle>各資產報酬率</CardTitle>
                <CardDescription>不同資產類別的投資績效比較</CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <BarChart data={mockPerformanceData}>
                    <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
                    <XAxis
                      dataKey="name"
                      className="text-xs"
                      tick={{ fill: 'hsl(var(--muted-foreground))' }}
                    />
                    <YAxis
                      className="text-xs"
                      tick={{ fill: 'hsl(var(--muted-foreground))' }}
                      label={{ value: '報酬率 (%)', angle: -90, position: 'insideLeft' }}
                    />
                    <Tooltip
                      contentStyle={{
                        backgroundColor: 'hsl(var(--background))',
                        border: '1px solid hsl(var(--border))',
                        borderRadius: '6px',
                      }}
                      formatter={(value: number) => [`${value.toFixed(2)}%`, '報酬率']}
                    />
                    <Bar dataKey="returnRate" radius={[4, 4, 0, 0]}>
                      {mockPerformanceData.map((entry, index) => (
                        <Cell
                          key={`cell-${index}`}
                          fill={entry.returnRate >= 0 ? '#10b981' : '#ef4444'}
                        />
                      ))}
                    </Bar>
                  </BarChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>

            {/* 資產配置圓餅圖 (重用 Dashboard 的資料) */}
            <Card>
              <CardHeader>
                <CardTitle>資產配置</CardTitle>
                <CardDescription>各資產類別的市值分布</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {mockAssetAllocation.map((item) => (
                    <div key={item.assetType} className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <div
                          className="h-3 w-3 rounded-full"
                          style={{ backgroundColor: item.color }}
                        />
                        <span className="text-sm font-medium">{item.name}</span>
                      </div>
                      <div className="flex items-center gap-4">
                        <span className="text-sm tabular-nums text-muted-foreground">
                          TWD {item.value.toLocaleString()}
                        </span>
                        <span className="text-sm font-medium tabular-nums w-12 text-right">
                          {item.percentage}%
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Top 資產表格 */}
          <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
            {/* Top 5 獲利資產 */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <TrendingUp className="h-5 w-5 text-green-600" />
                  Top 獲利資產
                </CardTitle>
                <CardDescription>表現最佳的投資標的</CardDescription>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>代碼/名稱</TableHead>
                      <TableHead>類別</TableHead>
                      <TableHead className="text-right">損益</TableHead>
                      <TableHead className="text-right">報酬率</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {mockTopProfitAssets.map((asset) => (
                      <TableRow key={asset.symbol}>
                        <TableCell>
                          <div>
                            <div className="font-medium">{asset.symbol}</div>
                            <div className="text-sm text-muted-foreground">{asset.name}</div>
                          </div>
                        </TableCell>
                        <TableCell>
                          <Badge variant="outline" className={assetTypeColors[asset.assetType]}>
                            {assetTypeNames[asset.assetType]}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-right font-medium tabular-nums text-green-600">
                          +{asset.profit.toLocaleString()}
                        </TableCell>
                        <TableCell className="text-right font-medium tabular-nums text-green-600">
                          +{asset.profitPercent}%
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>

            {/* Top 5 虧損資產 */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <TrendingDown className="h-5 w-5 text-red-600" />
                  Top 虧損資產
                </CardTitle>
                <CardDescription>需要關注的投資標的</CardDescription>
              </CardHeader>
              <CardContent>
                {mockTopLossAssets.length > 0 ? (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>代碼/名稱</TableHead>
                        <TableHead>類別</TableHead>
                        <TableHead className="text-right">損益</TableHead>
                        <TableHead className="text-right">報酬率</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {mockTopLossAssets.map((asset) => (
                        <TableRow key={asset.symbol}>
                          <TableCell>
                            <div>
                              <div className="font-medium">{asset.symbol}</div>
                              <div className="text-sm text-muted-foreground">{asset.name}</div>
                            </div>
                          </TableCell>
                          <TableCell>
                            <Badge variant="outline" className={assetTypeColors[asset.assetType]}>
                              {assetTypeNames[asset.assetType]}
                            </Badge>
                          </TableCell>
                          <TableCell className="text-right font-medium tabular-nums text-red-600">
                            {asset.profit.toLocaleString()}
                          </TableCell>
                          <TableCell className="text-right font-medium tabular-nums text-red-600">
                            {asset.profitPercent}%
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                ) : (
                  <div className="flex items-center justify-center h-32 text-muted-foreground">
                    目前沒有虧損資產
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        </div>
      </main>
    </AppLayout>
  );
}

