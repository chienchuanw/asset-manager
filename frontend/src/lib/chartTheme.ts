export const chartTheme = {
  gain: "var(--color-gain)",
  loss: "var(--color-loss)",
  neutral: "var(--color-neutral)",
  series: [
    "var(--color-chart-1)",
    "var(--color-chart-2)",
    "var(--color-chart-3)",
    "var(--color-chart-4)",
    "var(--color-chart-5)",
  ],
  grid: { stroke: "var(--color-surface-border)" },
  tooltip: {
    background: "var(--color-surface)",
    border: "var(--color-surface-border)",
    text: "var(--color-foreground)",
  },
} as const;
