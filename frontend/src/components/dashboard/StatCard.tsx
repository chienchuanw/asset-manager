/**
 * 統計卡片元件
 * 用於顯示總資產、今日損益等關鍵指標
 */

import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { TrendingUpIcon, TrendingDownIcon } from "lucide-react";

interface StatCardProps {
  title: string;
  value: string;
  change: number; // 百分比變化
  description?: string;
}

export function StatCard({ title, value, change, description }: StatCardProps) {
  // 判斷是上漲還是下跌
  const isPositive = change >= 0;

  return (
    <Card className="@container/card">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <p className="text-sm font-medium text-muted-foreground">{title}</p>
          <Badge
            variant="outline"
            className={`gap-1 ${
              isPositive
                ? "bg-red-50 text-red-700 border-red-200"
                : "bg-green-50 text-green-700 border-green-200"
            }`}
          >
            {isPositive ? (
              <TrendingUpIcon className="h-3 w-3" />
            ) : (
              <TrendingDownIcon className="h-3 w-3" />
            )}
            {isPositive ? "+" : ""}
            {change.toFixed(1)}%
          </Badge>
        </div>
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
          {value}
        </div>
      </CardContent>
      {description && (
        <CardFooter className="pt-0">
          <p className="text-xs text-muted-foreground">{description}</p>
        </CardFooter>
      )}
    </Card>
  );
}
