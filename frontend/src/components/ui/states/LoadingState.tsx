"use client";

import { Skeleton } from "@/components/ui/skeleton";

type Variant = "card" | "table" | "chart";

export function LoadingState({ variant = "card" }: { variant?: Variant }) {
  if (variant === "chart") {
    return <Skeleton className="h-[300px] w-full" />;
  }
  if (variant === "table") {
    return (
      <div className="space-y-2">
        {Array.from({ length: 5 }).map((_, i) => (
          <Skeleton key={i} className="h-10 w-full" />
        ))}
      </div>
    );
  }
  return (
    <div className="space-y-3">
      <Skeleton className="h-5 w-1/3" />
      <Skeleton className="h-8 w-2/3" />
    </div>
  );
}
