/**
 * 交易記錄頁面
 * 顯示所有交易記錄,支援篩選、搜尋、排序功能
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
import { Skeleton } from "@/components/ui/skeleton";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { AddTransactionDialog } from "@/components/transactions/AddTransactionDialog";
import { EditTransactionDialog } from "@/components/transactions/EditTransactionDialog";
import { useTransactions, useDeleteTransaction } from "@/hooks";
import {
  getAssetTypeLabel,
  getTransactionTypeLabel,
  AssetType,
  TransactionType,
  type Transaction,
} from "@/types/transaction";
import {
  Search,
  Download,
  Trash2,
  Loader2,
  MoreVertical,
  Edit,
} from "lucide-react";

export default function TransactionsPage() {
  // 狀態管理
  const [searchQuery, setSearchQuery] = useState("");
  const [filterType, setFilterType] = useState<TransactionType | "all">("all");
  const [filterAssetType, setFilterAssetType] = useState<AssetType | "all">(
    "all"
  );
  const [editingTransaction, setEditingTransaction] =
    useState<Transaction | null>(null);

  // 取得交易列表資料
  const { data: transactions, isLoading, error, refetch } = useTransactions();

  // 刪除交易 mutation
  const deleteMutation = useDeleteTransaction({
    onSuccess: () => {
      refetch();
    },
  });

  // 篩選和排序邏輯（使用 useMemo 優化效能）
  const filteredTransactions = useMemo(() => {
    if (!transactions) return [];

    return transactions
      .filter((transaction) => {
        // 交易類型篩選
        if (filterType !== "all" && transaction.type !== filterType) {
          return false;
        }
        // 資產類別篩選
        if (
          filterAssetType !== "all" &&
          transaction.asset_type !== filterAssetType
        ) {
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
        // 預設按日期降序排列（最新的在前）
        return new Date(b.date).getTime() - new Date(a.date).getTime();
      });
  }, [transactions, filterType, filterAssetType, searchQuery]);

  // 計算統計資料 (月度和日度)
  const stats = useMemo(() => {
    const now = new Date();
    const currentMonth = now.getMonth();
    const currentYear = now.getFullYear();
    const today = now.toDateString();

    // 本月交易
    const monthlyTransactions = filteredTransactions.filter((t) => {
      const txDate = new Date(t.date);
      return (
        txDate.getMonth() === currentMonth &&
        txDate.getFullYear() === currentYear
      );
    });

    // 今日交易
    const dailyTransactions = filteredTransactions.filter((t) => {
      const txDate = new Date(t.date);
      return txDate.toDateString() === today;
    });

    // 本月統計
    const monthlyBuyAmount = monthlyTransactions
      .filter((t) => t.type === TransactionType.BUY)
      .reduce((sum, t) => sum + t.amount, 0);
    const monthlySellAmount = monthlyTransactions
      .filter((t) => t.type === TransactionType.SELL)
      .reduce((sum, t) => sum + t.amount, 0);
    const monthlyNetFlow = monthlyBuyAmount - monthlySellAmount;

    // 今日統計
    const dailyBuyAmount = dailyTransactions
      .filter((t) => t.type === TransactionType.BUY)
      .reduce((sum, t) => sum + t.amount, 0);
    const dailySellAmount = dailyTransactions
      .filter((t) => t.type === TransactionType.SELL)
      .reduce((sum, t) => sum + t.amount, 0);
    const dailyNetFlow = dailyBuyAmount - dailySellAmount;

    return {
      monthlyTransactions: monthlyTransactions.length,
      monthlyBuyAmount,
      monthlySellAmount,
      monthlyNetFlow,
      dailyTransactions: dailyTransactions.length,
      dailyBuyAmount,
      dailySellAmount,
      dailyNetFlow,
    };
  }, [filteredTransactions]);

  // 處理刪除交易
  const handleDelete = (id: string) => {
    if (confirm("確定要刪除這筆交易嗎？")) {
      deleteMutation.mutate(id);
    }
  };

  // 取得交易類型的顏色
  const getTransactionTypeColor = (type: TransactionType) => {
    switch (type) {
      case TransactionType.BUY:
        return "bg-green-100 text-green-800";
      case TransactionType.SELL:
        return "bg-red-100 text-red-800";
      case TransactionType.DIVIDEND:
        return "bg-blue-100 text-blue-800";
      case TransactionType.FEE:
        return "bg-gray-100 text-gray-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  // 取得資產類型的顏色
  const getAssetTypeColor = (assetType: AssetType) => {
    switch (assetType) {
      case AssetType.TW_STOCK:
        return "bg-purple-100 text-purple-800";
      case AssetType.US_STOCK:
        return "bg-indigo-100 text-indigo-800";
      case AssetType.CRYPTO:
        return "bg-orange-100 text-orange-800";
      case AssetType.CASH:
        return "bg-emerald-100 text-emerald-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  return (
    <AppLayout title="交易記錄" description="管理和查看所有交易記錄">
      {/* Main Content */}
      <main className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="flex flex-col gap-6">
          {/* 統計摘要卡片 */}
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>本日交易次數</CardDescription>
              </CardHeader>
              <CardContent>
                {isLoading ? (
                  <Skeleton className="h-8 w-20" />
                ) : (
                  <>
                    <div className="text-2xl font-bold">
                      {stats.dailyTransactions}
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">
                      本月 {stats.monthlyTransactions} 筆
                    </p>
                  </>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardDescription>本日買入金額</CardDescription>
              </CardHeader>
              <CardContent>
                {isLoading ? (
                  <Skeleton className="h-8 w-32" />
                ) : (
                  <>
                    <div className="text-2xl font-bold text-blue-600">
                      TWD {stats.dailyBuyAmount.toLocaleString()}
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">
                      本月 TWD {stats.monthlyBuyAmount.toLocaleString()}
                    </p>
                  </>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardDescription>本日賣出金額</CardDescription>
              </CardHeader>
              <CardContent>
                {isLoading ? (
                  <Skeleton className="h-8 w-32" />
                ) : (
                  <>
                    <div className="text-2xl font-bold text-red-600">
                      TWD {stats.dailySellAmount.toLocaleString()}
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">
                      本月 TWD {stats.monthlySellAmount.toLocaleString()}
                    </p>
                  </>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardDescription>本日淨流入/流出</CardDescription>
              </CardHeader>
              <CardContent>
                {isLoading ? (
                  <Skeleton className="h-8 w-32" />
                ) : (
                  <>
                    <div
                      className={`text-2xl font-bold ${
                        stats.dailyNetFlow >= 0
                          ? "text-green-600"
                          : "text-red-600"
                      }`}
                    >
                      TWD {stats.dailyNetFlow.toLocaleString()}
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">
                      本月 TWD {stats.monthlyNetFlow.toLocaleString()}
                    </p>
                  </>
                )}
              </CardContent>
            </Card>
          </div>

          {/* 交易記錄表格 */}
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
                <Select
                  value={filterType}
                  onValueChange={(value) => setFilterType(value as any)}
                >
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
              {/* 錯誤訊息 */}
              {error && (
                <div className="p-4 mb-4 text-sm text-red-800 bg-red-100 rounded-lg">
                  <p className="font-medium">載入失敗</p>
                  <p>{error.message}</p>
                </div>
              )}

              <div className="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>日期</TableHead>
                      <TableHead>交易類型</TableHead>
                      <TableHead>資產類別</TableHead>
                      <TableHead>代碼/名稱</TableHead>
                      <TableHead className="text-right">數量</TableHead>
                      <TableHead className="text-right">單價</TableHead>
                      <TableHead className="text-right">總金額</TableHead>
                      <TableHead className="text-center">幣別</TableHead>
                      <TableHead className="text-right hidden md:table-cell">
                        手續費
                      </TableHead>
                      <TableHead className="hidden lg:table-cell">
                        備註
                      </TableHead>
                      <TableHead className="text-right">操作</TableHead>
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
                            <Skeleton className="h-6 w-16" />
                          </TableCell>
                          <TableCell>
                            <Skeleton className="h-6 w-16" />
                          </TableCell>
                          <TableCell>
                            <Skeleton className="h-4 w-32" />
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
                            <Skeleton className="h-6 w-12" />
                          </TableCell>
                          <TableCell className="hidden md:table-cell">
                            <Skeleton className="h-4 w-16" />
                          </TableCell>
                          <TableCell className="hidden lg:table-cell">
                            <Skeleton className="h-4 w-24" />
                          </TableCell>
                          <TableCell>
                            <Skeleton className="h-8 w-8" />
                          </TableCell>
                        </TableRow>
                      ))
                    ) : filteredTransactions.length === 0 ? (
                      // 無資料
                      <TableRow>
                        <TableCell colSpan={11} className="h-24 text-center">
                          <p className="text-muted-foreground">
                            {searchQuery ||
                            filterType !== "all" ||
                            filterAssetType !== "all"
                              ? "沒有符合條件的交易記錄"
                              : "尚無交易記錄，點擊「新增交易」開始記錄"}
                          </p>
                        </TableCell>
                      </TableRow>
                    ) : (
                      // 交易列表
                      filteredTransactions.map((transaction) => (
                        <TableRow key={transaction.id}>
                          <TableCell className="font-medium">
                            {new Date(transaction.date).toLocaleDateString(
                              "zh-TW"
                            )}
                          </TableCell>
                          <TableCell>
                            <Badge
                              variant="outline"
                              className={getTransactionTypeColor(
                                transaction.type
                              )}
                            >
                              {getTransactionTypeLabel(transaction.type)}
                            </Badge>
                          </TableCell>
                          <TableCell>
                            <Badge
                              variant="outline"
                              className={getAssetTypeColor(
                                transaction.asset_type
                              )}
                            >
                              {getAssetTypeLabel(transaction.asset_type)}
                            </Badge>
                          </TableCell>
                          <TableCell>
                            <div>
                              <div className="font-medium">
                                {transaction.symbol}
                              </div>
                              <div className="text-sm text-muted-foreground">
                                {transaction.name}
                              </div>
                            </div>
                          </TableCell>
                          <TableCell className="text-right">
                            {transaction.quantity.toLocaleString()}
                          </TableCell>
                          <TableCell className="text-right">
                            {transaction.price.toLocaleString()}
                          </TableCell>
                          <TableCell className="text-right font-medium">
                            {transaction.amount.toLocaleString()}
                          </TableCell>
                          <TableCell className="text-center">
                            <Badge
                              variant="outline"
                              className="bg-amber-100 text-amber-800"
                            >
                              {transaction.currency}
                            </Badge>
                          </TableCell>
                          <TableCell className="text-right hidden md:table-cell">
                            {transaction.fee
                              ? transaction.fee.toLocaleString()
                              : "-"}
                          </TableCell>
                          <TableCell className="text-sm text-muted-foreground hidden lg:table-cell">
                            {transaction.note || "-"}
                          </TableCell>
                          <TableCell className="text-right">
                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  disabled={deleteMutation.isPending}
                                >
                                  {deleteMutation.isPending ? (
                                    <Loader2 className="h-4 w-4 animate-spin" />
                                  ) : (
                                    <MoreVertical className="h-4 w-4" />
                                  )}
                                </Button>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align="end">
                                <DropdownMenuItem
                                  onClick={() => {
                                    setEditingTransaction(transaction);
                                  }}
                                >
                                  <Edit className="mr-2 h-4 w-4" />
                                  編輯
                                </DropdownMenuItem>
                                <DropdownMenuItem
                                  onClick={() => handleDelete(transaction.id)}
                                  className="text-red-600"
                                >
                                  <Trash2 className="mr-2 h-4 w-4" />
                                  刪除
                                </DropdownMenuItem>
                              </DropdownMenuContent>
                            </DropdownMenu>
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
