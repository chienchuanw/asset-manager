import { describe, it, expect } from "vitest";
import { chartTheme } from "./chartTheme";

describe("chartTheme", () => {
  it("exposes semantic series colors as CSS var references", () => {
    expect(chartTheme.gain).toBe("var(--color-gain)");
    expect(chartTheme.loss).toBe("var(--color-loss)");
  });
  it("exposes an ordered categorical series palette", () => {
    expect(Array.isArray(chartTheme.series)).toBe(true);
    expect(chartTheme.series.length).toBeGreaterThanOrEqual(5);
    expect(chartTheme.series[0]).toBe("var(--color-chart-1)");
  });
  it("exposes grid and tooltip styles bound to tokens", () => {
    expect(chartTheme.grid.stroke).toBe("var(--color-surface-border)");
    expect(chartTheme.tooltip.background).toBe("var(--color-surface)");
  });
});
