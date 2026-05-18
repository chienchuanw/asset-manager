import { describe, it, expect } from "vitest";
import { formatCurrency } from "./format";

describe("formatCurrency", () => {
  it("prefixes TWD with NT$", () => {
    expect(formatCurrency(1284500, "TWD")).toBe("NT$1,284,500");
  });
  it("prefixes USD with $", () => {
    expect(formatCurrency(1000, "USD")).toBe("$1,000");
  });
  it("defaults currency to TWD", () => {
    expect(formatCurrency(500)).toBe("NT$500");
  });
  it("renders negatives with sign before the prefix", () => {
    expect(formatCurrency(-2300, "TWD")).toBe("-NT$2,300");
  });
  it("rounds to whole units", () => {
    expect(formatCurrency(1234.56, "TWD")).toBe("NT$1,235");
  });
});
