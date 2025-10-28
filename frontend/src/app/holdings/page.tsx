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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Loading } from "@/components/ui/loading";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Search, ArrowUpDown, RefreshCw } from "lucide-react";
import { useHoldings } from "@/hooks";
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
  getConvertedStyle,
} from "@/types/holding";
import type { Holding } from "@/types/holding";

/**
 * 持倉卡片元件
 * 顯示單一資產類別的持倉列表
 */
interface HoldingCardProps {
  title: string;
  holdings: Holding[];
  showInTWD: boolean;
  sortConfig: {
    by: "market_value" | "unrealized_pl" | "quantity";
    order: "asc" | "desc";
  };
  onToggleSort: (field: "market_value" | "unrealized_pl" | "quantity") => void;
}

function HoldingCard({
  title,
  holdings,
  showInTWD,
  onToggleSort,
}: HoldingCardProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        <CardDescription>共 {holdings.length} 筆持倉</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>代碼/名稱</TableHead>
                <TableHead
                  className="cursor-pointer"
                  onClick={() => onToggleSort("quantity")}
                >
                  <div className="flex items-center gap-1">
                    數量
                    <ArrowUpDown className="h-3 w-3" />
                  </div>
                </TableHead>
                <TableHead
                  className="cursor-pointer text-right"
                  onClick={() => onToggleSort("market_value")}
                >
                  <div className="flex items-center justify-end gap-1">
                    市值
                    <ArrowUpDown className="h-3 w-3" />
                  </div>
                </TableHead>
                <TableHead
                  className="cursor-pointer text-right"
                  onClick={() => onToggleSort("unrealized_pl")}
                >
                  <div className="flex items-center justify-end gap-1">
                    損益
                    <ArrowUpDown className="h-3 w-3" />
                  </div>
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {holdings.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={4} className="h-24 text-center">
                    目前無持倉
                  </TableCell>
                </TableRow>
              ) : (
                holdings.map((holding) => (
                  <TableRow key={holding.symbol}>
                    <TableCell>
                      <div className="flex flex-col">
                        <span className="font-medium text-sm">
                          {holding.symbol}
                        </span>
                        <span className="text-xs text-muted-foreground">
                          {holding.name}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell className="tabular-nums text-sm">
                      {holding.quantity.toLocaleString("zh-TW", {
                        minimumFractionDigits: 0,
                        maximumFractionDigits: 4,
                      })}
                    </TableCell>
                    <TableCell className="text-right tabular-nums text-sm">
                      <span className={getConvertedStyle(showInTWD)}>
                        {showInTWD
                          ? formatCurrency(holding.market_value, "TWD")
                          : formatCurrency(
                              holding.market_value,
                              holding.currency
                            )}
                      </span>
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex flex-col items-end gap-0.5">
                        <span
                          className={`font-medium tabular-nums text-sm ${
                            showInTWD
                              ? getConvertedStyle(true, holding.unrealized_pl)
                              : getProfitLossColor(holding.unrealized_pl)
                          }`}
                        >
                          {showInTWD
                            ? formatCurrency(holding.unrealized_pl, "TWD")
                            : formatCurrency(
                                holding.unrealized_pl,
                                holding.currency
                              )}
                        </span>
                        <span
                          className={`text-xs tabular-nums ${
                            showInTWD
                              ? getConvertedStyle(
                                  true,
                                  holding.unrealized_pl_pct
                                )
                              : getProfitLossColor(holding.unrealized_pl_pct)
                          }`}
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
  );
}

export default function HoldingsPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [showInTWD, setShowInTWD] = useState(false);

  // 每個資產類別獨立的排序狀態
  const [twStockSort, setTwStockSort] = useState<{
    by: "market_value" | "unrealized_pl" | "quantity";
    order: "asc" | "desc";
  }>({ by: "market_value", order: "desc" });

  const [usStockSort, setUsStockSort] = useState<{
    by: "market_value" | "unrealized_pl" | "quantity";
    order: "asc" | "desc";
  }>({ by: "market_value", order: "desc" });

  const [cryptoSort, setCryptoSort] = useState<{
    by: "market_value" | "unrealized_pl" | "quantity";
    order: "asc" | "desc";
  }>({ by: "market_value", order: "desc" });

  // 從 API 取得所有持倉資料
  const {
    data: holdings,
    isLoading,
    error,
    refetch,
    isFetching,
  } = useHoldings();

  // 按資產類別分組並排序（使用 useMemo 優化效能）
  const holdingsByType = useMemo(() => {
    if (!holdings) return { twStock: [], usStock: [], crypto: [] };

    // 先搜尋
    const searched = searchHoldings(holdings, searchQuery);

    // 按資產類別分組
    const twStock = searched.filter((h) => h.asset_type === "tw-stock");
    const usStock = searched.filter((h) => h.asset_type === "us-stock");
    const crypto = searched.filter((h) => h.asset_type === "crypto");

    // 各自排序
    return {
      twStock: sortHoldings(twStock, twStockSort.by, twStockSort.order),
      usStock: sortHoldings(usStock, usStockSort.by, usStockSort.order),
      crypto: sortHoldings(crypto, cryptoSort.by, cryptoSort.order),
    };
  }, [holdings, searchQuery, twStockSort, usStockSort, cryptoSort]);

  // 計算統計資料（使用 useMemo 優化效能）
  const stats = useMemo(() => {
    if (!holdings || holdings.length === 0) {
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

    const totalMarketValue = calculateTotalMarketValue(holdings);
    const totalCost = calculateTotalCost(holdings);
    const totalProfitLoss = calculateTotalProfitLoss(holdings);
    const totalProfitLossPercent = calculateTotalProfitLossPct(holdings);

    // 計算各類資產市值
    const twStockValue = holdingsByType.twStock.reduce(
      (sum, h) => sum + h.market_value,
      0
    );
    const usStockValue = holdingsByType.usStock.reduce(
      (sum, h) => sum + h.market_value,
      0
    );
    const cryptoValue = holdingsByType.crypto.reduce(
      (sum, h) => sum + h.market_value,
      0
    );

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
  }, [holdings, holdingsByType]);

  // 切換排序函式（每個資產類別獨立）
  const toggleTwStockSort = (
    field: "market_value" | "unrealized_pl" | "quantity"
  ) => {
    setTwStockSort((prev) => ({
      by: field,
      order: prev.by === field && prev.order === "desc" ? "asc" : "desc",
    }));
  };

  const toggleUsStockSort = (
    field: "market_value" | "unrealized_pl" | "quantity"
  ) => {
    setUsStockSort((prev) => ({
      by: field,
      order: prev.by === field && prev.order === "desc" ? "asc" : "desc",
    }));
  };

  const toggleCryptoSort = (
    field: "market_value" | "unrealized_pl" | "quantity"
  ) => {
    setCryptoSort((prev) => ({
      by: field,
      order: prev.by === field && prev.order === "desc" ? "asc" : "desc",
    }));
  };

  // Loading 狀態
  if (isLoading) {
    return (
      <AppLayout title="持倉明細" description="查看所有資產的詳細持倉資訊">
        <div className="flex-1 p-4 md:p-6 bg-gray-50">
          <div className="flex items-center justify-center h-96">
            <Loading variant="page" size="lg" text="載入持倉資料中..." />
          </div>
        </div>
      </AppLayout>
    );
  }

  // 錯誤狀態
  if (error) {
    return (
      <AppLayout title="持倉明細" description="查看所有資產的詳細持倉資訊">
        <div className="flex-1 p-4 md:p-6 bg-gray-50">
          <div className="flex items-center justify-center h-96">
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
        </div>
      </AppLayout>
    );
  }

  return (
    <AppLayout title="持倉明細" description="查看所有資產的詳細持倉資訊">
      {/* Main Content */}
      <div className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="flex flex-col gap-6">
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

          {/* 搜尋與控制工具列 */}
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
              {/* 搜尋框 */}
              <div className="relative">
                <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="搜尋代碼或名稱..."
                  className="pl-8 sm:w-[250px]"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                />
              </div>
            </div>

            <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
              {/* 重新整理按鈕 */}
              <Button
                variant="outline"
                size="sm"
                onClick={() => refetch()}
                disabled={isFetching}
              >
                <RefreshCw
                  className={`mr-2 h-4 w-4 ${isFetching ? "animate-spin" : ""}`}
                />
                重新整理
              </Button>

              {/* 幣別切換開關 */}
              <div className="flex items-center gap-2">
                <Switch
                  id="currency-toggle"
                  checked={showInTWD}
                  onCheckedChange={setShowInTWD}
                />
                <Label
                  htmlFor="currency-toggle"
                  className="text-sm cursor-pointer"
                >
                  TWD
                </Label>
              </div>
            </div>
          </div>
          {/* 三張持倉卡片 */}
          <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
            {/* 台股持倉卡片 */}
            <HoldingCard
              title="台股持倉"
              holdings={holdingsByType.twStock}
              showInTWD={showInTWD}
              sortConfig={twStockSort}
              onToggleSort={toggleTwStockSort}
            />

            {/* 美股持倉卡片 */}
            <HoldingCard
              title="美股持倉"
              holdings={holdingsByType.usStock}
              showInTWD={showInTWD}
              sortConfig={usStockSort}
              onToggleSort={toggleUsStockSort}
            />

            {/* 加密貨幣持倉卡片 */}
            <HoldingCard
              title="加密貨幣持倉"
              holdings={holdingsByType.crypto}
              showInTWD={showInTWD}
              sortConfig={cryptoSort}
              onToggleSort={toggleCryptoSort}
            />
          </div>
        </div>
      </div>
    </AppLayout>
  );
}
