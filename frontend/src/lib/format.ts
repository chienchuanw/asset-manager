export function formatCurrency(value: number, currency: string = "TWD"): string {
  const symbol = currency === "TWD" ? "NT$" : currency === "USD" ? "$" : null;
  const sign = value < 0 ? "-" : "";
  const abs = Math.abs(value).toLocaleString("zh-TW", {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  });
  // Known currencies get a glyph prefix; anything else keeps its code so the
  // amount is never rendered without a currency indicator.
  return symbol ? `${sign}${symbol}${abs}` : `${sign}${currency} ${abs}`;
}
