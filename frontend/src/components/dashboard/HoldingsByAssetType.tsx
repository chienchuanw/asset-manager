/**
 * 依資產類別分組的持倉卡片元件
 * 將持倉分成台股、美股、加密貨幣三張獨立卡片
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
import { Holding, getProfitLossColor } from '@/types/holding';
import { AssetType } from '@/types/transaction';

interface HoldingsByAssetTypeProps {
  holdings: Holding[];
}

interface AssetTypeCardProps {
  title: string;
  holdings: Holding[];
  emptyMessage: string;
}

function AssetTypeCard({ title, holdings, emptyMessage }: AssetTypeCardProps) {
  // 計算總計
  const totals = holdings.reduce(
    (acc, holding) => ({
      marketValue: acc.marketValue + holding.market_value,
      unrealizedPL: acc.unrealizedPL + holding.unrealized_pl,
    }),
    { marketValue: 0, unrealizedPL: 0 }
  );

  const totalPLPct = totals.marketValue > 0 
    ? ((totals.unrealizedPL / (totals.marketValue - totals.unrealizedPL)) * 100)
    : 0;

  const profitLossColor = getProfitLossColor(totals.unrealizedPL);

  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        <CardDescription>
          {holdings.length > 0 ? (
            <div className="flex items-center gap-4 mt-2">
              <span>共 {holdings.length} 筆持倉</span>
              <span className="text-muted-foreground">|</span>
              <span>市值: TWD {totals.marketValue.toLocaleString('zh-TW', { maximumFractionDigits: 0 })}</span>
              <span className="text-muted-foreground">|</span>
              <span className={profitLossColor}>
                損益: {totals.unrealizedPL >= 0 ? '+' : ''}
                {totals.unrealizedPL.toLocaleString('zh-TW', { maximumFractionDigits: 0 })}
                ({totals.unrealizedPL >= 0 ? '+' : ''}{totalPLPct.toFixed(2)}%)
              </span>
            </div>
          ) : (
            emptyMessage
          )}
        </CardDescription>
      </CardHeader>
      <CardContent>
        {holdings.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            {emptyMessage}
          </div>
        ) : (
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>資產</TableHead>
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
        )}
      </CardContent>
    </Card>
  );
}

export function HoldingsByAssetType({ holdings }: HoldingsByAssetTypeProps) {
  // 依資產類別分組
  const twStockHoldings = holdings.filter(h => h.asset_type === AssetType.TW_STOCK);
  const usStockHoldings = holdings.filter(h => h.asset_type === AssetType.US_STOCK);
  const cryptoHoldings = holdings.filter(h => h.asset_type === AssetType.CRYPTO);

  return (
    <div className="space-y-4">
      <AssetTypeCard
        title="台股持倉"
        holdings={twStockHoldings}
        emptyMessage="目前無台股持倉"
      />
      <AssetTypeCard
        title="美股持倉"
        holdings={usStockHoldings}
        emptyMessage="目前無美股持倉"
      />
      <AssetTypeCard
        title="加密貨幣持倉"
        holdings={cryptoHoldings}
        emptyMessage="目前無加密貨幣持倉"
      />
    </div>
  );
}

