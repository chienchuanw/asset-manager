/**
 * 資產配置偏差表格元件
 * 顯示各資產類別的當前配置、目標配置和偏差
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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type { AssetTypeDeviation } from "@/types/rebalance";
import { ASSET_TYPE_MAP } from "./constants";

interface AllocationTableProps {
  deviations: AssetTypeDeviation[];
  threshold: number;
}

export function AllocationTable({
  deviations,
  threshold,
}: AllocationTableProps) {
  const t = useTranslations("rebalance");
  const tAssets = useTranslations("assetTypes");
  const tCommon = useTranslations("common");

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
   * 格式化百分比
   */
  function formatPercent(value: number): string {
    const sign = value > 0 ? "+" : "";
    return `${sign}${value.toFixed(2)}%`;
  }

  /**
   * 格式化金額
   */
  function formatAmount(value: number): string {
    return `NT$ ${value.toLocaleString("zh-TW")}`;
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("currentAllocation")}</CardTitle>
        <CardDescription>
          {t("threshold")}: {threshold}%
        </CardDescription>
      </CardHeader>
      <CardContent>
        {deviations.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            {tCommon("noData")}
          </div>
        ) : (
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{t("assetType")}</TableHead>
                  <TableHead className="text-right">
                    {t("currentAllocation")}
                  </TableHead>
                  <TableHead className="text-right">
                    {t("targetAllocation")}
                  </TableHead>
                  <TableHead className="text-right">
                    {t("deviation")}
                  </TableHead>
                  <TableHead className="text-right">
                    {t("amount")}
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {deviations.map((deviation) => (
                  <TableRow key={deviation.asset_type}>
                    <TableCell className="font-medium">
                      {getAssetTypeName(deviation.asset_type)}
                    </TableCell>
                    <TableCell className="text-right tabular-nums">
                      {deviation.current_percent.toFixed(2)}%
                    </TableCell>
                    <TableCell className="text-right tabular-nums">
                      {deviation.target_percent.toFixed(2)}%
                    </TableCell>
                    <TableCell className="text-right">
                      <span
                        data-exceeds-threshold={String(
                          deviation.exceeds_threshold
                        )}
                        className={`tabular-nums font-medium ${
                          deviation.exceeds_threshold
                            ? "text-red-600"
                            : "text-muted-foreground"
                        }`}
                      >
                        {formatPercent(deviation.deviation)}
                      </span>
                    </TableCell>
                    <TableCell className="text-right tabular-nums text-sm text-muted-foreground">
                      {formatAmount(deviation.current_value)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
