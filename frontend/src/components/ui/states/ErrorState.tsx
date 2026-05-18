"use client";

import { useTranslations } from "next-intl";
import { AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";

export function ErrorState({
  message,
  onRetry,
}: {
  message?: string;
  onRetry?: () => void;
}) {
  const t = useTranslations("states");
  return (
    <div className="flex h-full min-h-[200px] w-full flex-col items-center justify-center gap-3 text-center">
      <AlertCircle className="h-8 w-8 text-loss" />
      <p className="text-sm">{message ?? t("errorTitle")}</p>
      {onRetry && (
        <Button variant="outline" size="sm" onClick={onRetry}>
          {t("retry")}
        </Button>
      )}
    </div>
  );
}
