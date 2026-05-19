/**
 * FIRE 計算頁面
 * 顯示 FIRE 數字、達標進度、預估達標年數，以及淨值與 FI 目標投影圖。
 */

"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { AppLayout } from "@/components/layout/AppLayout";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { Loading } from "@/components/ui/loading";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Money } from "@/components/common/Money";
import {
  Area,
  ComposedChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { AlertCircle } from "lucide-react";
import { useFireProjection } from "@/hooks/useFireProjection";
import type { FireProjectionParams, FireProjectionResult } from "@/types/fire";

type DraftKey = keyof FireProjectionParams;

const DRAFT_FIELDS: DraftKey[] = [
  "netWorth",
  "annualExpenses",
  "annualSavings",
  "expectedReturn",
  "withdrawalRate",
];

export default function FirePage() {
  const t = useTranslations("fire");
  const tCommon = useTranslations("common");
  const tErrors = useTranslations("errors");

  const [params, setParams] = useState<FireProjectionParams>({});
  const [draft, setDraft] = useState<Record<DraftKey, string>>({
    netWorth: "",
    annualExpenses: "",
    annualSavings: "",
    expectedReturn: "",
    withdrawalRate: "",
  });

  const { data, isLoading, isError } = useFireProjection(params);

  const handleApply = () => {
    const next: FireProjectionParams = {};
    for (const key of DRAFT_FIELDS) {
      const raw = draft[key].trim();
      if (raw === "") continue;
      const parsed = Number(raw);
      if (!Number.isNaN(parsed)) next[key] = parsed;
    }
    setParams(next);
  };

  const handleReset = () => {
    setDraft({
      netWorth: "",
      annualExpenses: "",
      annualSavings: "",
      expectedReturn: "",
      withdrawalRate: "",
    });
    setParams({});
  };

  const progressValue =
    data && data.fire_number > 0
      ? Math.min(100, Math.max(0, data.progress_pct * 100))
      : data
        ? 100
        : 0;

  return (
    <AppLayout title={t("title")} description={t("description")}>
      <div className="flex-1 p-4 md:p-6 bg-muted">
        <div className="flex flex-col gap-6">
          {/* 假設輸入表單 */}
        <Card>
          <CardHeader>
            <CardTitle>{t("assumptions")}</CardTitle>
            <CardDescription>{t("assumptionsHint")}</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-5">
              {DRAFT_FIELDS.map((key) => (
                <div key={key} className="space-y-2">
                  <Label htmlFor={key}>{t(`fields.${key}`)}</Label>
                  <Input
                    id={key}
                    inputMode="decimal"
                    placeholder={
                      data ? String(echoedValue(data, key)) : tCommon("loading")
                    }
                    value={draft[key]}
                    onChange={(e) =>
                      setDraft((d) => ({ ...d, [key]: e.target.value }))
                    }
                  />
                </div>
              ))}
            </div>
            <div className="mt-4 flex gap-2">
              <Button onClick={handleApply}>{t("apply")}</Button>
              <Button variant="outline" onClick={handleReset}>
                {t("resetToAuto")}
              </Button>
            </div>
          </CardContent>
        </Card>

        {isLoading && <Loading />}

        {isError && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>{tErrors("loadFailed")}</AlertTitle>
            <AlertDescription>{tErrors("serverError")}</AlertDescription>
          </Alert>
        )}

        {data && !isLoading && (
          <>
            {data.current_net_worth === 0 && (
              <Alert>
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>{t("noSnapshotTitle")}</AlertTitle>
                <AlertDescription>{t("noSnapshotHint")}</AlertDescription>
              </Alert>
            )}

            {/* 結果卡片 */}
            <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
              <Card>
                <CardHeader>
                  <CardDescription>{t("fireNumber")}</CardDescription>
                  <CardTitle className="text-2xl">
                    <Money value={data.fire_number} />
                  </CardTitle>
                </CardHeader>
              </Card>
              <Card>
                <CardHeader>
                  <CardDescription>{t("progress")}</CardDescription>
                  <CardTitle className="text-2xl">
                    {(data.progress_pct * 100).toFixed(1)}%
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <Progress value={progressValue} />
                </CardContent>
              </Card>
              <Card>
                <CardHeader>
                  <CardDescription>{t("yearsToFI")}</CardDescription>
                  <CardTitle className="text-2xl">
                    {data.on_track && data.years_to_fi !== null
                      ? t("yearsValue", { years: data.years_to_fi })
                      : t("notOnTrack")}
                  </CardTitle>
                </CardHeader>
              </Card>
            </div>

            {/* 投影圖 */}
            <Card>
              <CardHeader>
                <CardTitle>{t("projectionTitle")}</CardTitle>
                <CardDescription>{t("projectionHint")}</CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={320}>
                  <ComposedChart data={data.projection}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis
                      dataKey="year"
                      tickFormatter={(y) => t("yearAxis", { year: y })}
                    />
                    <YAxis
                      width={80}
                      tickFormatter={(v) =>
                        new Intl.NumberFormat("en", {
                          notation: "compact",
                        }).format(v as number)
                      }
                    />
                    <Tooltip
                      formatter={(value: number) =>
                        new Intl.NumberFormat("en").format(value)
                      }
                    />
                    <Area
                      type="monotone"
                      dataKey="projected_net_worth"
                      name={t("fields.netWorth")}
                      fill="hsl(var(--primary) / 0.2)"
                      stroke="hsl(var(--primary))"
                    />
                    <Line
                      type="monotone"
                      dataKey="fire_target"
                      name={t("fireNumber")}
                      stroke="hsl(var(--destructive))"
                      dot={false}
                      strokeDasharray="5 5"
                    />
                  </ComposedChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>
          </>
        )}
        </div>
      </div>
    </AppLayout>
  );
}

function echoedValue(data: FireProjectionResult, key: DraftKey): number {
  const map: Record<DraftKey, keyof FireProjectionResult> = {
    netWorth: "current_net_worth",
    annualExpenses: "annual_expenses",
    annualSavings: "annual_savings",
    expectedReturn: "expected_return",
    withdrawalRate: "withdrawal_rate",
  };
  const v = data[map[key]];
  return typeof v === "number" ? v : 0;
}
