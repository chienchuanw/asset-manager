/**
 * 再平衡建議卡片元件
 * 顯示買入/賣出建議或已平衡訊息
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
import { CheckCircle2Icon } from "lucide-react";
import type { RebalanceSuggestion } from "@/types/rebalance";

/**
 * 資產類別顯示名稱對應
 */
const ASSET_TYPE_MAP: Record<string, string> = {
  "tw-stock": "twStock",
  "us-stock": "usStock",
  crypto: "crypto",
};

interface SuggestionCardsProps {
  suggestions: RebalanceSuggestion[];
}

export function SuggestionCards({ suggestions }: SuggestionCardsProps) {
  const t = useTranslations("rebalance");
  const tAssets = useTranslations("assetTypes");

  /**
   * 取得資產類別的顯示名稱
   */
  function getAssetTypeName(assetType: string): string {
    const key = ASSET_TYPE_MAP[assetType];
    if (key) {
      return tAssets(key as "twStock" | "usStock" | "crypto");
    }
    return assetType;
  }

  /**
   * 格式化金額
   */
  function formatAmount(value: number): string {
    return `NT$ ${value.toLocaleString("zh-TW")}`;
  }

  // 已平衡，無需調整
  if (suggestions.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>{t("suggestions")}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center justify-center py-8 gap-3 text-center">
            <CheckCircle2Icon className="h-12 w-12 text-green-500" />
            <p className="text-muted-foreground">{t("balanced")}</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("suggestions")}</CardTitle>
        <CardDescription>{t("description")}</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="flex flex-col gap-3">
          {suggestions.map((suggestion, index) => (
            <div
              key={`${suggestion.asset_type}-${index}`}
              className="flex items-start gap-4 rounded-lg border p-4"
            >
              <Badge
                variant={suggestion.action === "buy" ? "default" : "destructive"}
                className="mt-0.5"
              >
                {suggestion.action === "buy" ? t("buy") : t("sell")}
              </Badge>
              <div className="flex-1 space-y-1">
                <div className="flex items-center justify-between">
                  <span className="font-medium">
                    {getAssetTypeName(suggestion.asset_type)}
                  </span>
                  <span className="font-medium tabular-nums">
                    {formatAmount(suggestion.amount)}
                  </span>
                </div>
                <p className="text-sm text-muted-foreground">
                  {suggestion.reason}
                </p>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
