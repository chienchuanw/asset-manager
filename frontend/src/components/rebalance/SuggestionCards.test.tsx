import { screen } from "@testing-library/react";
import { renderWithProviders } from "@/test/utils";
import { SuggestionCards } from "./SuggestionCards";
import type { RebalanceSuggestion } from "@/types/rebalance";

const mockSuggestions: RebalanceSuggestion[] = [
  {
    asset_type: "tw-stock",
    action: "sell",
    amount: 100000,
    reason: "台股配置超過目標 10%，建議賣出",
  },
  {
    asset_type: "us-stock",
    action: "buy",
    amount: 50000,
    reason: "美股配置低於目標 5%，建議買入",
  },
];

describe("SuggestionCards", () => {
  it("should render buy and sell suggestion cards", () => {
    renderWithProviders(<SuggestionCards suggestions={mockSuggestions} />);

    // 檢查建議行動
    expect(screen.getByText("Sell")).toBeInTheDocument();
    expect(screen.getByText("Buy")).toBeInTheDocument();

    // 檢查資產類別
    expect(screen.getByText("TW Stock")).toBeInTheDocument();
    expect(screen.getByText("US Stock")).toBeInTheDocument();

    // 檢查原因
    expect(
      screen.getByText("台股配置超過目標 10%，建議賣出")
    ).toBeInTheDocument();
    expect(
      screen.getByText("美股配置低於目標 5%，建議買入")
    ).toBeInTheDocument();
  });

  it("should show balanced message when no suggestions", () => {
    renderWithProviders(<SuggestionCards suggestions={[]} />);

    expect(screen.getByText("Your portfolio is well balanced. No rebalancing is needed at this time.")).toBeInTheDocument();
  });

  it("should display amounts formatted as TWD", () => {
    renderWithProviders(<SuggestionCards suggestions={mockSuggestions} />);

    // 金額應有千分位
    expect(screen.getByText("NT$ 100,000")).toBeInTheDocument();
    expect(screen.getByText("NT$ 50,000")).toBeInTheDocument();
  });
});
