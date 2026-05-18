"use client";

import { useTranslations } from "next-intl";
import { Inbox } from "lucide-react";

export function EmptyState({
  title,
  description,
  action,
}: {
  title?: string;
  description?: string;
  action?: React.ReactNode;
}) {
  const t = useTranslations("states");
  return (
    <div className="flex h-full min-h-[200px] w-full flex-col items-center justify-center gap-2 text-center">
      <Inbox className="h-8 w-8 text-muted-foreground" />
      <p className="text-sm font-medium">{title ?? t("emptyTitle")}</p>
      <p className="text-xs text-muted-foreground">
        {description ?? t("emptyDescription")}
      </p>
      {action}
    </div>
  );
}
