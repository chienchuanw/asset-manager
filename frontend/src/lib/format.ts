export function formatCurrency(value: number, currency: string = "TWD"): string {
  const prefix = currency === "TWD" ? "NT$" : currency === "USD" ? "$" : "";
  const sign = value < 0 ? "-" : "";
  const abs = Math.abs(value).toLocaleString("zh-TW", {
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  });
  return `${sign}${prefix}${abs}`;
}
