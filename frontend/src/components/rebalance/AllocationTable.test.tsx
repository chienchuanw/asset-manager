import { screen } from "@testing-library/react";
import { renderWithProviders } from "@/test/utils";
import { AllocationTable } from "./AllocationTable";
import type { AssetTypeDeviation } from "@/types/rebalance";

const mockDeviations: AssetTypeDeviation[] = [
  {
    asset_type: "tw-stock",
    target_percent: 50,
    current_percent: 60,
    deviation: 10,
    deviation_abs: 10,
    exceeds_threshold: true,
    current_value: 600000,
    target_value: 500000,
  },
  {
    asset_type: "us-stock",
    target_percent: 30,
    current_percent: 28,
    deviation: -2,
    deviation_abs: 2,
    exceeds_threshold: false,
    current_value: 280000,
    target_value: 300000,
  },
  {
    asset_type: "crypto",
    target_percent: 20,
    current_percent: 12,
    deviation: -8,
    deviation_abs: 8,
    exceeds_threshold: true,
    current_value: 120000,
    target_value: 200000,
  },
];

describe("AllocationTable", () => {
  it("should render all asset types with percentages", () => {
    renderWithProviders(
      <AllocationTable deviations={mockDeviations} threshold={5} />
    );

    // 檢查資產類別是否顯示
    expect(screen.getByText("TW Stock")).toBeInTheDocument();
    expect(screen.getByText("US Stock")).toBeInTheDocument();
    expect(screen.getByText("Crypto")).toBeInTheDocument();

    // 檢查百分比是否顯示
    expect(screen.getByText("60.00%")).toBeInTheDocument();
    expect(screen.getByText("50.00%")).toBeInTheDocument();
  });

  it("should highlight asset types that exceed threshold", () => {
    renderWithProviders(
      <AllocationTable deviations={mockDeviations} threshold={5} />
    );

    // 超過閾值的偏差應該有特殊樣式（紅色文字）
    const twStockDeviation = screen.getByText("+10.00%");
    expect(twStockDeviation).toBeInTheDocument();
    expect(twStockDeviation.closest("[data-exceeds-threshold]")).toHaveAttribute(
      "data-exceeds-threshold",
      "true"
    );

    // 未超過閾值的偏差
    const usStockDeviation = screen.getByText("-2.00%");
    expect(usStockDeviation).toBeInTheDocument();
    expect(usStockDeviation.closest("[data-exceeds-threshold]")).toHaveAttribute(
      "data-exceeds-threshold",
      "false"
    );
  });

  it("should render empty state when no deviations", () => {
    renderWithProviders(
      <AllocationTable deviations={[]} threshold={5} />
    );

    expect(screen.getByText("No data")).toBeInTheDocument();
  });
});
