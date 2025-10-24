/**
 * 持倉明細頁面
 * 顯示所有資產的詳細持倉資訊,支援篩選、排序、搜尋功能
 */

"use client";

import { useState, useMemo } from "react";
import { AppLayout } from "@/components/layout/AppLayout";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Search, ArrowUpDown, Download, RefreshCw } from "lucide-react";
import { useHoldings } from "@/hooks";
import { AssetType } from "@/types/transaction";
import { getAssetTypeLabel } from "@/types/transaction";
import {
  sortHoldings,
  searchHoldings,
  calculateTotalMarketValue,
  calculateTotalCost,
  calculateTotalProfitLoss,
  calculateTotalProfitLossPct,
  formatCurrency,
  formatPercentage,
  getProfitLossColor,
} from "@/types/holding";

export default function HoldingsPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [filterType, setFilterType] = useState<AssetType | "all">("all");
  const [sortBy, setSortBy] = useState<
    "market_value" | "unrealized_pl" | "quantity"
  >("market_value");
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">("desc");

  // 從 API 取得持倉資料
  const {
    data: holdings,
    isLoading,
    error,
    refetch,
    isFetching,
  } = useHoldings(
    filterType !== "all" ? { asset_type: filterType } : undefined
  );

  // 篩選與排序邏輯（使用 useMemo 優化效能）
  const filteredAndSortedHoldings = useMemo(() => {
    if (!holdings) return [];

    // 先搜尋
    const searched = searchHoldings(holdings, searchQuery);

    // 再排序
    return sortHoldings(searched, sortBy, sortOrder);
  }, [holdings, searchQuery, sortBy, sortOrder]);

  // 計算統計資料（使用 useMemo 優化效能）
  const stats = useMemo(() => {
    if (!filteredAndSortedHoldings.length) {
      return {
        totalMarketValue: 0,
        totalCost: 0,
        totalProfitLoss: 0,
        totalProfitLossPercent: 0,
        twStockValue: 0,
        usStockValue: 0,
        cryptoValue: 0,
        availableCash: 0,
      };
    }

    const totalMarketValue = calculateTotalMarketValue(
      filteredAndSortedHoldings
    );
    const totalCost = calculateTotalCost(filteredAndSortedHoldings);
    const totalProfitLoss = calculateTotalProfitLoss(filteredAndSortedHoldings);
    const totalProfitLossPercent = calculateTotalProfitLossPct(
      filteredAndSortedHoldings
    );

    // 計算各類資產市值
    const twStockValue = filteredAndSortedHoldings
      .filter((h) => h.asset_type === "tw-stock")
      .reduce((sum, h) => sum + h.market_value, 0);

    const usStockValue = filteredAndSortedHoldings
      .filter((h) => h.asset_type === "us-stock")
      .reduce((sum, h) => sum + h.market_value, 0);

    const cryptoValue = filteredAndSortedHoldings
      .filter((h) => h.asset_type === "crypto")
      .reduce((sum, h) => sum + h.market_value, 0);

    // 計算可用現金 (總市值 - 總成本 = 未實現損益)
    const availableCash = totalProfitLoss;

    return {
      totalMarketValue,
      totalCost,
      totalProfitLoss,
      totalProfitLossPercent,
      twStockValue,
      usStockValue,
      cryptoValue,
      availableCash,
    };
  }, [filteredAndSortedHoldings]);

  // 切換排序
  const toggleSort = (field: "market_value" | "unrealized_pl" | "quantity") => {
    if (sortBy === field) {
      setSortOrder(sortOrder === "asc" ? "desc" : "asc");
    } else {
      setSortBy(field);
      setSortOrder("desc");
    }
  };

  // Loading 狀態
  if (isLoading) {
    return (
      <AppLayout>
        <main className="flex-1 p-4 md:p-6 bg-gray-50">
          <div className="container flex items-center justify-center h-96">
            <div className="flex flex-col items-center gap-4">
              <RefreshCw className="h-8 w-8 animate-spin text-muted-foreground" />
              <p className="text-muted-foreground">載入持倉資料中...</p>
            </div>
          </div>
        </main>
      </AppLayout>
    );
  }

  // 錯誤狀態
  if (error) {
    return (
      <AppLayout>
        <main className="flex-1 p-4 md:p-6 bg-gray-50">
          <div className="container flex items-center justify-center h-96">
            <Card className="w-full max-w-md">
              <CardHeader>
                <CardTitle className="text-red-600">載入失敗</CardTitle>
                <CardDescription>
                  無法載入持倉資料：{error.message}
                </CardDescription>
              </CardHeader>
              <CardContent>
                <Button
                  onClick={() => refetch()}
                  variant="outline"
                  className="w-full"
                >
                  <RefreshCw className="mr-2 h-4 w-4" />
                  重新載入
                </Button>
              </CardContent>
            </Card>
          </div>
        </main>
      </AppLayout>
    );
  }

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
                  NT$ {stats.totalMarketValue.toLocaleString()}
                </CardTitle>
              </CardHeader>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>總成本</CardDescription>
                <CardTitle className="text-2xl tabular-nums">
                  NT$ {stats.totalCost.toLocaleString()}
                </CardTitle>
              </CardHeader>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>未實現損益</CardDescription>
                <CardTitle
                  className={`text-2xl tabular-nums ${getProfitLossColor(
                    stats.totalProfitLoss
                  )}`}
                >
                  {formatCurrency(stats.totalProfitLoss, "TWD")}
                </CardTitle>
              </CardHeader>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>報酬率</CardDescription>
                <CardTitle
                  className={`text-2xl tabular-nums ${getProfitLossColor(
                    stats.totalProfitLossPercent
                  )}`}
                >
                  {formatPercentage(stats.totalProfitLossPercent)}
                </CardTitle>
              </CardHeader>
            </Card>
          </div>

          {/* 統計摘要卡片 - 第二列 */}
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>台股市值</CardDescription>
                <CardTitle className="text-2xl tabular-nums">
                  NT$ {stats.twStockValue.toLocaleString()}
                </CardTitle>
              </CardHeader>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>美股市值</CardDescription>
                <CardTitle className="text-2xl tabular-nums">
                  NT$ {stats.usStockValue.toLocaleString()}
                </CardTitle>
              </CardHeader>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>加密貨幣市值</CardDescription>
                <CardTitle className="text-2xl tabular-nums">
                  NT$ {stats.cryptoValue.toLocaleString()}
                </CardTitle>
              </CardHeader>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>可用現金</CardDescription>
                <CardTitle
                  className={`text-2xl tabular-nums ${getProfitLossColor(
                    stats.availableCash
                  )}`}
                >
                  {formatCurrency(stats.availableCash, "TWD")}
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
                    {isFetching && (
                      <span className="ml-2 text-xs">(更新中...)</span>
                    )}
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
                    onValueChange={(value) =>
                      setFilterType(value as AssetType | "all")
                    }
                  >
                    <SelectTrigger className="sm:w-[150px]">
                      <SelectValue placeholder="資產類別" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">全部類別</SelectItem>
                      <SelectItem value="tw-stock">台股</SelectItem>
                      <SelectItem value="us-stock">美股</SelectItem>
                      <SelectItem value="crypto">加密貨幣</SelectItem>
                    </SelectContent>
                  </Select>

                  {/* 重新整理按鈕 */}
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => refetch()}
                    disabled={isFetching}
                  >
                    <RefreshCw
                      className={`mr-2 h-4 w-4 ${
                        isFetching ? "animate-spin" : ""
                      }`}
                    />
                    重新整理
                  </Button>

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
                        onClick={() => toggleSort("quantity")}
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
                        onClick={() => toggleSort("market_value")}
                      >
                        <div className="flex items-center justify-end gap-1">
                          市值
                          <ArrowUpDown className="h-3 w-3" />
                        </div>
                      </TableHead>
                      <TableHead
                        className="cursor-pointer text-right"
                        onClick={() => toggleSort("unrealized_pl")}
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
                        <TableRow key={holding.symbol}>
                          <TableCell>
                            <div className="flex flex-col">
                              <span className="font-medium">
                                {holding.symbol}
                              </span>
                              <span className="text-sm text-muted-foreground">
                                {holding.name}
                              </span>
                            </div>
                          </TableCell>
                          <TableCell>
                            <Badge variant="outline">
                              {getAssetTypeLabel(holding.asset_type)}
                            </Badge>
                          </TableCell>
                          <TableCell className="tabular-nums">
                            {holding.quantity.toLocaleString("zh-TW", {
                              minimumFractionDigits: 0,
                              maximumFractionDigits: 8,
                            })}
                          </TableCell>
                          <TableCell className="text-right tabular-nums">
                            {holding.avg_cost.toLocaleString("zh-TW", {
                              minimumFractionDigits: 0,
                              maximumFractionDigits: 2,
                            })}
                          </TableCell>
                          <TableCell className="text-right tabular-nums">
                            {holding.current_price.toLocaleString("zh-TW", {
                              minimumFractionDigits: 0,
                              maximumFractionDigits: 2,
                            })}
                          </TableCell>
                          <TableCell className="text-right tabular-nums">
                            {formatCurrency(
                              holding.market_value,
                              holding.currency
                            )}
                          </TableCell>
                          <TableCell className="text-right">
                            <div className="flex flex-col items-end">
                              <span
                                className={`font-medium tabular-nums ${getProfitLossColor(
                                  holding.unrealized_pl
                                )}`}
                              >
                                {formatCurrency(
                                  holding.unrealized_pl,
                                  holding.currency
                                )}
                              </span>
                              <span
                                className={`text-sm tabular-nums ${getProfitLossColor(
                                  holding.unrealized_pl_pct
                                )}`}
                              >
                                {formatPercentage(holding.unrealized_pl_pct)}
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
