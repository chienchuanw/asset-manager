/**
 * 交易記錄頁面
 * 顯示交易記錄,支援篩選、搜尋、排序功能
 * 優化手機介面,預設顯示今日記錄
 */

"use client";

import { useState, useMemo } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { AppLayout } from "@/components/layout/AppLayout";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { AddTransactionDialog } from "@/components/transactions/AddTransactionDialog";
import { BatchAddTransactionDialog } from "@/components/transactions/BatchAddTransactionDialog";
import { EditTransactionDialog } from "@/components/transactions/EditTransactionDialog";
import { CSVImportDialog } from "@/components/transactions/CSVImportDialog";
import { DailyDateNavigator } from "@/components/cash-flows/DailyDateNavigator";
import { WeekMonthTabs } from "@/components/common/WeekMonthTabs";
import { calculateDateRange } from "@/components/common/DateRangeTabs";
import {
  useTransactions,
  useDeleteTransaction,
  transactionKeys,
} from "@/hooks";
import {
  TransactionType,
  type Transaction,
  type TransactionFilters,
  getAssetTypeLabel,
} from "@/types/transaction";
import { Search, MoreVertical, Pencil, Trash2 } from "lucide-react";

export default function TransactionsPage() {
  const queryClient = useQueryClient();

  // 狀態管理
  const [searchQuery, setSearchQuery] = useState("");

  // 上半部：Tab 控制統計卡片的日期範圍
  const [statsTab, setStatsTab] = useState<"week" | "month">("week");

  // 下半部：日期導航控制交易列表顯示的日期
  const [selectedDate, setSelectedDate] = useState<Date>(new Date());

  const [editingTransaction, setEditingTransaction] =
    useState<Transaction | null>(null);
  const [showBatchDialog, setShowBatchDialog] = useState(false);
  const [csvTransactions, setCsvTransactions] = useState<any[]>([]);

  // 計算上半部統計卡片的日期範圍
  const { startDate: statsStartDate, endDate: statsEndDate } = useMemo(() => {
    return calculateDateRange(statsTab);
  }, [statsTab]);

  // 計算下半部交易列表的日期範圍（只顯示選定的那一天）
  const { startDate: listStartDate, endDate: listEndDate } = useMemo(() => {
    // 使用本地時間格式化日期，避免時區轉換問題
    const year = selectedDate.getFullYear();
    const month = String(selectedDate.getMonth() + 1).padStart(2, "0");
    const day = String(selectedDate.getDate()).padStart(2, "0");
    const dateStr = `${year}-${month}-${day}`;
    return {
      startDate: dateStr,
      endDate: dateStr,
    };
  }, [selectedDate]);

  // 建立統計卡片的 API 篩選條件
  const statsFilters: TransactionFilters = useMemo(() => {
    return {
      start_date: statsStartDate,
      end_date: statsEndDate,
    };
  }, [statsStartDate, statsEndDate]);

  // 建立交易列表的 API 篩選條件
  const listFilters: TransactionFilters = useMemo(() => {
    return {
      start_date: listStartDate,
      end_date: listEndDate,
    };
  }, [listStartDate, listEndDate]);

  // 取得統計卡片的交易資料
  const { data: statsTransactions } = useTransactions(statsFilters, {
    staleTime: 0,
  });

  // 取得交易列表的交易資料
  const {
    data: transactions,
    isLoading,
    error,
  } = useTransactions(listFilters, {
    staleTime: 0,
  });

  // 刪除交易 mutation
  const deleteMutation = useDeleteTransaction({
    onSuccess: () => {
      handleRefreshData();
    },
  });

  // 客戶端搜尋篩選（只針對交易列表）
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

  // 計算統計資料（基於統計卡片的日期範圍）
  const stats = useMemo(() => {
    if (!statsTransactions) {
      return {
        totalTransactions: 0,
        buyAmount: 0,
        sellAmount: 0,
        netFlow: 0,
      };
    }

    const buyAmount = statsTransactions
      .filter((t) => t.type === TransactionType.BUY)
      .reduce((sum, t) => sum + t.amount, 0);

    const sellAmount = statsTransactions
      .filter((t) => t.type === TransactionType.SELL)
      .reduce((sum, t) => sum + t.amount, 0);

    const netFlow = buyAmount - sellAmount;

    return {
      totalTransactions: statsTransactions.length,
      buyAmount,
      sellAmount,
      netFlow,
    };
  }, [statsTransactions]);

  // 處理刪除交易
  const handleDelete = (id: string) => {
    if (confirm("確定要刪除這筆交易嗎？")) {
      deleteMutation.mutate(id);
    }
  };

  // 重新獲取所有相關資料
  const handleRefreshData = async () => {
    // 讓所有交易相關查詢失效，強制重新獲取
    await queryClient.invalidateQueries({
      queryKey: transactionKeys.all,
    });
  };

  // 處理 CSV 匯入成功
  const handleCSVImportSuccess = (transactions: any[]) => {
    setCsvTransactions(transactions);
    setShowBatchDialog(true);
  };

  return (
    <AppLayout title="交易記錄" description="管理和查看交易記錄">
      {/* Main Content */}
      <div className="flex-1 p-4 md:p-6 bg-gray-50 space-y-6">
        {/* 上半部：Tab + 統計卡片 */}
        <div className="flex flex-col gap-6">
          {/* Tab 切換 */}
          <WeekMonthTabs value={statsTab} onValueChange={setStatsTab} />

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
        </div>

        {/* 下半部：日期導航 + 交易記錄列表 */}
        <div className="flex flex-col gap-6">
          {/* 日期導航 */}
          <DailyDateNavigator
            date={selectedDate}
            onDateChange={setSelectedDate}
          />

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
                  <AddTransactionDialog onSuccess={handleRefreshData} />
                  <BatchAddTransactionDialog onSuccess={handleRefreshData} />
                  <CSVImportDialog onSuccess={handleCSVImportSuccess} />
                </div>
              </div>

              {/* 搜尋框 */}
              <div className="mt-4">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                  <Input
                    placeholder="搜尋代碼或名稱..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="pl-9"
                  />
                </div>
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

              {/* 交易記錄表格 */}
              <div className="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>日期</TableHead>
                      <TableHead>資產</TableHead>
                      <TableHead>類型</TableHead>
                      <TableHead className="text-right">數量</TableHead>
                      <TableHead className="text-right">價格</TableHead>
                      <TableHead className="text-right">金額</TableHead>
                      <TableHead className="text-right hidden md:table-cell">
                        手續費
                      </TableHead>
                      <TableHead className="text-right hidden md:table-cell">
                        交易稅
                      </TableHead>
                      <TableHead className="w-[50px]"></TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {isLoading ? (
                      // 載入中骨架屏
                      Array.from({ length: 5 }).map((_, index) => (
                        <TableRow key={index}>
                          <TableCell>
                            <Skeleton className="h-4 w-20" />
                          </TableCell>
                          <TableCell>
                            <Skeleton className="h-4 w-24" />
                          </TableCell>
                          <TableCell>
                            <Skeleton className="h-4 w-16" />
                          </TableCell>
                          <TableCell>
                            <Skeleton className="h-4 w-16" />
                          </TableCell>
                          <TableCell>
                            <Skeleton className="h-4 w-20" />
                          </TableCell>
                          <TableCell>
                            <Skeleton className="h-4 w-20" />
                          </TableCell>
                          <TableCell className="hidden md:table-cell">
                            <Skeleton className="h-4 w-16" />
                          </TableCell>
                          <TableCell className="hidden md:table-cell">
                            <Skeleton className="h-4 w-16" />
                          </TableCell>
                          <TableCell>
                            <Skeleton className="h-4 w-8" />
                          </TableCell>
                        </TableRow>
                      ))
                    ) : filteredTransactions.length === 0 ? (
                      // 無資料
                      <TableRow>
                        <TableCell colSpan={9} className="h-24 text-center">
                          <p className="text-muted-foreground">
                            {searchQuery
                              ? "沒有符合條件的交易記錄"
                              : "尚無交易記錄，點擊「新增交易」開始記錄"}
                          </p>
                        </TableCell>
                      </TableRow>
                    ) : (
                      // 交易記錄列表
                      filteredTransactions.map((transaction) => {
                        const typeColor =
                          transaction.type === "buy"
                            ? "text-red-600"
                            : transaction.type === "sell"
                            ? "text-green-600"
                            : transaction.type === "dividend"
                            ? "text-blue-600"
                            : "text-gray-600";

                        const typeLabel =
                          transaction.type === "buy"
                            ? "買入"
                            : transaction.type === "sell"
                            ? "賣出"
                            : transaction.type === "dividend"
                            ? "股利"
                            : "手續費";

                        return (
                          <TableRow key={transaction.id}>
                            <TableCell className="text-sm">
                              {new Date(transaction.date).toLocaleDateString(
                                "zh-TW"
                              )}
                            </TableCell>
                            <TableCell>
                              <div className="flex flex-col">
                                <span className="font-medium text-sm">
                                  {transaction.symbol}
                                </span>
                                <span className="text-xs text-muted-foreground">
                                  {transaction.name}
                                </span>
                                <Badge
                                  variant="outline"
                                  className="w-fit mt-1 text-xs"
                                >
                                  {getAssetTypeLabel(transaction.asset_type)}
                                </Badge>
                              </div>
                            </TableCell>
                            <TableCell>
                              <Badge
                                variant="outline"
                                className={`${typeColor} border-current`}
                              >
                                {typeLabel}
                              </Badge>
                            </TableCell>
                            <TableCell className="text-right tabular-nums text-sm">
                              {transaction.quantity?.toLocaleString("zh-TW", {
                                minimumFractionDigits: 0,
                                maximumFractionDigits: 4,
                              }) || "-"}
                            </TableCell>
                            <TableCell className="text-right tabular-nums text-sm">
                              {transaction.price
                                ? `${transaction.price.toLocaleString("zh-TW", {
                                    minimumFractionDigits: 2,
                                    maximumFractionDigits: 4,
                                  })} ${transaction.currency}`
                                : "-"}
                            </TableCell>
                            <TableCell className="text-right tabular-nums text-sm font-medium">
                              {transaction.amount.toLocaleString("zh-TW", {
                                minimumFractionDigits: 2,
                              })}{" "}
                              <span className="text-muted-foreground font-normal">
                                {transaction.currency}
                              </span>
                            </TableCell>
                            <TableCell className="text-right tabular-nums text-sm text-muted-foreground hidden md:table-cell">
                              {transaction.fee
                                ? `${transaction.fee.toLocaleString("zh-TW", {
                                    minimumFractionDigits: 2,
                                  })} ${transaction.currency}`
                                : "-"}
                            </TableCell>
                            <TableCell className="text-right tabular-nums text-sm text-muted-foreground hidden md:table-cell">
                              {transaction.tax
                                ? `${transaction.tax.toLocaleString("zh-TW", {
                                    minimumFractionDigits: 2,
                                  })} ${transaction.currency}`
                                : "-"}
                            </TableCell>
                            <TableCell>
                              <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                  <Button variant="ghost" size="sm">
                                    <MoreVertical className="h-4 w-4" />
                                  </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                  <DropdownMenuItem
                                    onClick={() =>
                                      setEditingTransaction(transaction)
                                    }
                                  >
                                    <Pencil className="h-4 w-4 mr-2" />
                                    編輯
                                  </DropdownMenuItem>
                                  <DropdownMenuItem
                                    onClick={() => handleDelete(transaction.id)}
                                    className="text-red-600"
                                    disabled={deleteMutation.isPending}
                                  >
                                    <Trash2 className="h-4 w-4 mr-2" />
                                    刪除
                                  </DropdownMenuItem>
                                </DropdownMenuContent>
                              </DropdownMenu>
                            </TableCell>
                          </TableRow>
                        );
                      })
                    )}
                  </TableBody>
                </Table>
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
          onSuccess={handleRefreshData}
        />
      )}

      {/* CSV 匯入後的批量新增對話框 */}
      {showBatchDialog && (
        <BatchAddTransactionDialog
          open={showBatchDialog}
          onOpenChange={setShowBatchDialog}
          onSuccess={handleRefreshData}
          initialTransactions={csvTransactions}
        />
      )}
    </AppLayout>
  );
}
