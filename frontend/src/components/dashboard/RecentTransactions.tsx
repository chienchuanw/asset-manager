/**
 * 近期交易列表元件
 * 顯示最近的交易記錄
 */

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Transaction, assetTypeNames } from '@/lib/mock-data';
import { ArrowUpIcon, ArrowDownIcon } from 'lucide-react';

interface RecentTransactionsProps {
  transactions: Transaction[];
}

export function RecentTransactions({ transactions }: RecentTransactionsProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>近期交易</CardTitle>
        <CardDescription>最近的買賣記錄</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {transactions.map((transaction) => {
            const isBuy = transaction.type === 'buy';
            return (
              <div
                key={transaction.id}
                className="flex items-start justify-between pb-4 border-b last:border-0 last:pb-0"
              >
                <div className="flex items-start gap-3">
                  {/* 買賣圖示 */}
                  <div
                    className={`mt-0.5 p-1 rounded-full ${
                      isBuy ? 'bg-green-100' : 'bg-red-100'
                    }`}
                  >
                    {isBuy ? (
                      <ArrowDownIcon className="h-3 w-3 text-green-600" />
                    ) : (
                      <ArrowUpIcon className="h-3 w-3 text-red-600" />
                    )}
                  </div>

                  {/* 交易資訊 */}
                  <div className="flex-1 space-y-1">
                    <div className="flex items-center gap-2">
                      <span className="font-medium text-sm">
                        {transaction.symbol}
                      </span>
                      <Badge
                        variant="outline"
                        className="text-xs"
                      >
                        {isBuy ? '買入' : '賣出'}
                      </Badge>
                    </div>
                    <p className="text-sm text-muted-foreground">{transaction.name}</p>
                    <p className="text-xs text-muted-foreground">
                      {transaction.date}
                    </p>
                  </div>
                </div>

                {/* 金額 */}
                <div className="text-right space-y-1">
                  <p className="font-medium text-sm tabular-nums">
                    {transaction.amount.toLocaleString()}
                  </p>
                  <p className="text-xs text-muted-foreground tabular-nums">
                    {transaction.quantity} × {transaction.price.toLocaleString()}
                  </p>
                </div>
              </div>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}

