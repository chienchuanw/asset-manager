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
  it("keeps up to two fraction digits so currency cents survive", () => {
    expect(formatCurrency(1234.56, "TWD")).toBe("NT$1,234.56");
    expect(formatCurrency(1000, "USD")).toBe("$1,000");
  });
  it("falls back to the currency code for unmapped currencies", () => {
    expect(formatCurrency(1000, "EUR")).toBe("EUR 1,000");
    expect(formatCurrency(-50.5, "JPY")).toBe("-JPY 50.5");
  });
});
