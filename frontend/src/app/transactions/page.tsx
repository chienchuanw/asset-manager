/**
 * 交易記錄頁面
 * 顯示交易記錄,支援篩選、搜尋、排序功能
 * 優化手機介面,預設顯示今日記錄
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
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { AddTransactionDialog } from "@/components/transactions/AddTransactionDialog";
import { BatchAddTransactionDialog } from "@/components/transactions/BatchAddTransactionDialog";
import { EditTransactionDialog } from "@/components/transactions/EditTransactionDialog";
import { TransactionCard } from "@/components/transactions/TransactionCard";
import { TransactionFilterDrawer } from "@/components/transactions/TransactionFilterDrawer";
import {
  DateRangeTabs,
  calculateDateRange,
  type DateRangeType,
} from "@/components/common/DateRangeTabs";
import { useTransactions, useDeleteTransaction } from "@/hooks";
import {
  AssetType,
  TransactionType,
  type Transaction,
  type TransactionFilters,
} from "@/types/transaction";
import { Search, Download } from "lucide-react";
import { DateRange } from "react-day-picker";

export default function TransactionsPage() {
  // 狀態管理
  const [searchQuery, setSearchQuery] = useState("");
  const [dateRangeType, setDateRangeType] = useState<DateRangeType>("today");
  const [customDateRange, setCustomDateRange] = useState<DateRange | undefined>(
    undefined
  );
  const [filterType, setFilterType] = useState<TransactionType | "all">("all");
  const [filterAssetType, setFilterAssetType] = useState<AssetType | "all">(
    "all"
  );
  const [editingTransaction, setEditingTransaction] =
    useState<Transaction | null>(null);

  // 計算日期範圍
  const { startDate, endDate } = useMemo(() => {
    // 如果有自訂日期範圍,優先使用
    if (customDateRange?.from) {
      return {
        startDate: customDateRange.from.toISOString().split("T")[0],
        endDate: customDateRange.to
          ? customDateRange.to.toISOString().split("T")[0]
          : customDateRange.from.toISOString().split("T")[0],
      };
    }
    // 否則使用 Tabs 選擇的日期範圍
    return calculateDateRange(dateRangeType);
  }, [dateRangeType, customDateRange]);

  // 建立 API 篩選條件
  const apiFilters: TransactionFilters = useMemo(() => {
    const filters: TransactionFilters = {
      start_date: startDate,
      end_date: endDate,
    };

    if (filterType !== "all") {
      filters.type = filterType;
    }

    if (filterAssetType !== "all") {
      filters.asset_type = filterAssetType;
    }

    return filters;
  }, [startDate, endDate, filterType, filterAssetType]);

  // 取得交易列表資料
  const {
    data: transactions,
    isLoading,
    error,
    refetch,
  } = useTransactions(apiFilters);

  // 刪除交易 mutation
  const deleteMutation = useDeleteTransaction({
    onSuccess: () => {
      refetch();
    },
  });

  // 客戶端搜尋篩選（API 已經處理了日期和類型篩選）
  const filteredTransactions = useMemo(() => {
    if (!transactions) return [];

    return transactions
      .filter((transaction) => {
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
        // 預設按日期降序排列（最新的在前）
        return new Date(b.date).getTime() - new Date(a.date).getTime();
      });
  }, [transactions, searchQuery]);

  // 計算統計資料（基於當前篩選的交易）
  const stats = useMemo(() => {
    if (!filteredTransactions) {
      return {
        totalTransactions: 0,
        buyAmount: 0,
        sellAmount: 0,
        netFlow: 0,
      };
    }

    const buyAmount = filteredTransactions
      .filter((t) => t.type === TransactionType.BUY)
      .reduce((sum, t) => sum + t.amount, 0);

    const sellAmount = filteredTransactions
      .filter((t) => t.type === TransactionType.SELL)
      .reduce((sum, t) => sum + t.amount, 0);

    const netFlow = buyAmount - sellAmount;

    return {
      totalTransactions: filteredTransactions.length,
      buyAmount,
      sellAmount,
      netFlow,
    };
  }, [filteredTransactions]);

  // 處理刪除交易
  const handleDelete = (id: string) => {
    if (confirm("確定要刪除這筆交易嗎？")) {
      deleteMutation.mutate(id);
    }
  };

  // 處理重置篩選
  const handleResetFilters = () => {
    setFilterType("all");
    setFilterAssetType("all");
    setCustomDateRange(undefined);
  };

  return (
    <AppLayout title="交易記錄" description="管理和查看交易記錄">
      {/* Main Content */}
      <div className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="flex flex-col gap-6">
          {/* 日期範圍 Tabs */}
          <DateRangeTabs
            value={dateRangeType}
            onValueChange={(value) => {
              setDateRangeType(value);
              setCustomDateRange(undefined); // 切換 Tabs 時清除自訂日期
            }}
          />

          {/* 統計摘要卡片 */}
          <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>交易次數</CardDescription>
              </CardHeader>
              <CardContent>
                {isLoading ? (
                  <Skeleton className="h-8 w-20" />
                ) : (
                  <div className="text-2xl font-bold">
                    {stats.totalTransactions}
                  </div>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardDescription>買入金額</CardDescription>
              </CardHeader>
              <CardContent>
                {isLoading ? (
                  <Skeleton className="h-8 w-32" />
                ) : (
                  <div className="text-2xl font-bold text-blue-600">
                    TWD {stats.buyAmount.toLocaleString()}
                  </div>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardDescription>賣出金額</CardDescription>
              </CardHeader>
              <CardContent>
                {isLoading ? (
                  <Skeleton className="h-8 w-32" />
                ) : (
                  <div className="text-2xl font-bold text-red-600">
                    TWD {stats.sellAmount.toLocaleString()}
                  </div>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardDescription>淨流入/流出</CardDescription>
              </CardHeader>
              <CardContent>
                {isLoading ? (
                  <Skeleton className="h-8 w-32" />
                ) : (
                  <div
                    className={`text-2xl font-bold ${
                      stats.netFlow >= 0 ? "text-green-600" : "text-red-600"
                    }`}
                  >
                    TWD {stats.netFlow.toLocaleString()}
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          {/* 交易記錄列表 */}
          <Card>
            <CardHeader>
              <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <CardTitle>交易記錄</CardTitle>
                  <CardDescription>
                    {isLoading
                      ? "載入中..."
                      : `共 ${filteredTransactions.length} 筆記錄`}
                  </CardDescription>
                </div>
                <div className="flex gap-2">
                  <AddTransactionDialog onSuccess={() => refetch()} />
                  <BatchAddTransactionDialog onSuccess={() => refetch()} />
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

                {/* 進階篩選 Drawer */}
                <TransactionFilterDrawer
                  filterType={filterType}
                  filterAssetType={filterAssetType}
                  dateRange={customDateRange}
                  onFilterTypeChange={setFilterType}
                  onFilterAssetTypeChange={setFilterAssetType}
                  onDateRangeChange={setCustomDateRange}
                  onReset={handleResetFilters}
                />
              </div>
            </CardHeader>

            <CardContent>
              {/* 錯誤訊息 */}
              {error && (
                <div className="p-4 mb-4 text-sm text-red-800 bg-red-100 rounded-lg">
                  <p className="font-medium">載入失敗</p>
                  <p>{error.message}</p>
                </div>
              )}

              {/* 交易記錄卡片列表 */}
              <div className="space-y-4">
                {isLoading ? (
                  // 載入中骨架屏
                  Array.from({ length: 3 }).map((_, index) => (
                    <Card key={index}>
                      <CardContent className="p-4">
                        <Skeleton className="h-32 w-full" />
                      </CardContent>
                    </Card>
                  ))
                ) : filteredTransactions.length === 0 ? (
                  // 無資料
                  <div className="text-center py-12">
                    <p className="text-muted-foreground">
                      {searchQuery ||
                      filterType !== "all" ||
                      filterAssetType !== "all" ||
                      customDateRange
                        ? "沒有符合條件的交易記錄"
                        : "尚無交易記錄，點擊「新增交易」開始記錄"}
                    </p>
                  </div>
                ) : (
                  // 交易卡片列表
                  filteredTransactions.map((transaction) => (
                    <TransactionCard
                      key={transaction.id}
                      transaction={transaction}
                      onEdit={setEditingTransaction}
                      onDelete={handleDelete}
                      isDeleting={deleteMutation.isPending}
                    />
                  ))
                )}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* 編輯交易對話框 */}
      {editingTransaction && (
        <EditTransactionDialog
          transaction={editingTransaction}
          open={!!editingTransaction}
          onOpenChange={(open) => {
            if (!open) setEditingTransaction(null);
          }}
          onSuccess={() => refetch()}
        />
      )}
    </AppLayout>
  );
}
