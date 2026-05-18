"use client";

import { useZenMode } from "@/providers/ZenModeProvider";
import { formatCurrency } from "@/lib/format";

type MoneyProps = {
  value: number;
  currency?: string;
  className?: string;
};

export function Money({ value, currency = "TWD", className }: MoneyProps) {
  const { zen } = useZenMode();
  return (
    <span className={className}>
      {zen ? "NT$•••••" : formatCurrency(value, currency)}
    </span>
  );
}
