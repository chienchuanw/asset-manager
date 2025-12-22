/**
 * 近期交易列表元件
 * 顯示最近的交易記錄
 */

"use client";

import { useTranslations } from "next-intl";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Transaction, AssetType, TransactionType } from "@/types/transaction";
import { ArrowUpIcon, ArrowDownIcon } from "lucide-react";

interface RecentTransactionsProps {
  transactions: Transaction[];
}

export function RecentTransactions({ transactions }: RecentTransactionsProps) {
  const t = useTranslations("dashboard");
  const tTx = useTranslations("transactions");
  const tAssets = useTranslations("assetTypes");

  // 取得資產類型的翻譯標籤
  const getAssetLabel = (assetType: string): string => {
    const keyMap: Record<string, string> = {
      "tw-stock": "twStock",
      "us-stock": "usStock",
      crypto: "crypto",
      cash: "cash",
    };
    const key = keyMap[assetType] || assetType;
    return tAssets(key as keyof typeof tAssets);
  };

  // 取得交易類型的翻譯標籤
  const getTransactionLabel = (type: string): string => {
    const keyMap: Record<string, string> = {
      buy: "buy",
      sell: "sell",
      dividend: "dividend",
      fee: "fee",
    };
    const key = keyMap[type] || type;
    return tTx(key as keyof typeof tTx);
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("recentTransactions")}</CardTitle>
        <CardDescription>{t("recentTransactionsDesc")}</CardDescription>
      </CardHeader>
      <CardContent>
        {transactions.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            {tTx("noTransactions")}
          </div>
        ) : (
          <div className="space-y-4">
            {transactions.map((transaction) => {
              const isBuy = transaction.type === "buy";
              return (
                <div
                  key={transaction.id}
                  className="flex items-start justify-between pb-4 border-b last:border-0 last:pb-0"
                >
                  <div className="flex items-start gap-3">
                    {/* 買賣圖示 */}
                    <div
                      className={`mt-0.5 p-1 rounded-full ${
                        isBuy ? "bg-red-100" : "bg-green-100"
                      }`}
                    >
                      {isBuy ? (
                        <ArrowDownIcon className="h-3 w-3 text-red-600" />
                      ) : (
                        <ArrowUpIcon className="h-3 w-3 text-green-600" />
                      )}
                    </div>

                    {/* 交易資訊 */}
                    <div className="flex-1 space-y-1">
                      <div className="flex items-center gap-2">
                        <span className="font-medium text-sm">
                          {transaction.symbol}
                        </span>
                        <Badge variant="outline" className="text-xs">
                          {getTransactionLabel(transaction.type)}
                        </Badge>
                      </div>
                      <p className="text-sm text-muted-foreground">
                        {getAssetLabel(transaction.asset_type)}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        {new Date(transaction.date).toLocaleDateString("zh-TW")}
                      </p>
                    </div>
                  </div>

                  {/* 金額 */}
                  <div className="text-right space-y-1">
                    <p className="font-medium text-sm tabular-nums">
                      {transaction.currency}{" "}
                      {transaction.amount.toLocaleString("zh-TW", {
                        maximumFractionDigits: 0,
                      })}
                    </p>
                    <p className="text-xs text-muted-foreground tabular-nums">
                      {transaction.quantity.toLocaleString("zh-TW", {
                        maximumFractionDigits: 2,
                      })}{" "}
                      ×{" "}
                      {transaction.price.toLocaleString("zh-TW", {
                        maximumFractionDigits: 2,
                      })}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
