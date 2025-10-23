/**
 * 持倉明細表格元件
 * 顯示所有持倉的詳細資訊
 */

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Holding } from '@/types/holding';
import { getAssetTypeLabel, getProfitLossColor } from '@/types/transaction';

interface HoldingsTableProps {
  holdings: Holding[];
}

export function HoldingsTable({ holdings }: HoldingsTableProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>持倉明細</CardTitle>
        <CardDescription>目前持有的所有資產</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>資產</TableHead>
                <TableHead className="hidden sm:table-cell">類別</TableHead>
                <TableHead className="text-right">數量</TableHead>
                <TableHead className="text-right hidden md:table-cell">成本價</TableHead>
                <TableHead className="text-right">現價</TableHead>
                <TableHead className="text-right">市值</TableHead>
                <TableHead className="text-right">損益</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {holdings.map((holding) => {
                const isProfit = holding.unrealized_pl >= 0;
                const profitLossColor = getProfitLossColor(holding.unrealized_pl);
                return (
                  <TableRow key={holding.symbol}>
                    <TableCell>
                      <div className="font-medium">{holding.symbol}</div>
                      <div className="text-sm text-muted-foreground hidden md:block">
                        {holding.name}
                      </div>
                    </TableCell>
                    <TableCell className="hidden sm:table-cell">
                      <Badge variant="outline" className="text-xs">
                        {getAssetTypeLabel(holding.asset_type)}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right tabular-nums">
                      {holding.quantity.toLocaleString('zh-TW', { maximumFractionDigits: 2 })}
                    </TableCell>
                    <TableCell className="text-right tabular-nums hidden md:table-cell">
                      {holding.avg_cost.toLocaleString('zh-TW', { maximumFractionDigits: 2 })}
                    </TableCell>
                    <TableCell className="text-right tabular-nums">
                      {holding.current_price_twd.toLocaleString('zh-TW', { maximumFractionDigits: 2 })}
                    </TableCell>
                    <TableCell className="text-right font-medium tabular-nums">
                      {holding.market_value.toLocaleString('zh-TW', { maximumFractionDigits: 0 })}
                    </TableCell>
                    <TableCell className={`text-right font-medium tabular-nums ${profitLossColor}`}>
                      <div>{isProfit ? '+' : ''}{holding.unrealized_pl.toLocaleString('zh-TW', { maximumFractionDigits: 0 })}</div>
                      <div className="text-xs">
                        {isProfit ? '+' : ''}{holding.unrealized_pl_pct.toFixed(2)}%
                      </div>
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  );
}

