// Mock window.matchMedia (required by shadcn Sidebar / use-mobile hook)
Object.defineProperty(window, "matchMedia", {
  writable: true,
  value: vi.fn().mockImplementation((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});

import { screen, waitFor } from "@testing-library/react";
import { http, HttpResponse } from "msw";
import { server } from "@/test/server";
import { renderWithProviders } from "@/test/utils";
import RebalancePage from "./page";
import type { RebalanceCheck } from "@/types/rebalance";

const API_BASE = "http://localhost:8080";

// Mock Next.js navigation
vi.mock("next/navigation", () => ({
  usePathname: () => "/rebalance",
  useRouter: () => ({ push: vi.fn(), replace: vi.fn(), back: vi.fn() }),
}));

// Mock next/image
vi.mock("next/image", () => ({
  default: (props: Record<string, unknown>) => {
    const { src, alt, ...rest } = props;
    return <img src={src as string} alt={alt as string} {...rest} />;
  },
}));

// Mock AuthProvider
vi.mock("@/providers/AuthProvider", () => ({
  useAuth: () => ({ logout: vi.fn() }),
}));

// Mock LocaleProvider (used by LanguageSwitcher in AppLayout)
vi.mock("@/providers/LocaleProvider", () => ({
  useLocale: () => ({ locale: "en", setLocale: vi.fn() }),
}));

const mockRebalanceData: RebalanceCheck = {
  needs_rebalance: true,
  threshold: 5,
  deviations: [
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
      current_percent: 25,
      deviation: -5,
      deviation_abs: 5,
      exceeds_threshold: true,
      current_value: 250000,
      target_value: 300000,
    },
  ],
  suggestions: [
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
  ],
  current_total: 1000000,
};

const balancedData: RebalanceCheck = {
  needs_rebalance: false,
  threshold: 5,
  deviations: [
    {
      asset_type: "tw-stock",
      target_percent: 50,
      current_percent: 51,
      deviation: 1,
      deviation_abs: 1,
      exceeds_threshold: false,
      current_value: 510000,
      target_value: 500000,
    },
  ],
  suggestions: [],
  current_total: 1000000,
};

describe("RebalancePage", () => {
  it("should show loading state initially", () => {
    server.use(
      http.get(`${API_BASE}/api/rebalance/check`, () => {
        return new Promise(() => {});
      })
    );

    renderWithProviders(<RebalancePage />);

    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  it("should show allocation table and suggestions when needs_rebalance is true", async () => {
    server.use(
      http.get(`${API_BASE}/api/rebalance/check`, () => {
        return HttpResponse.json({ data: mockRebalanceData });
      })
    );

    renderWithProviders(<RebalancePage />);

    // 等待資料載入完成
    await waitFor(() => {
      expect(screen.getByText("NT$ 1,000,000")).toBeInTheDocument();
    });

    // 摘要卡片
    expect(screen.getByText("Total Portfolio Value")).toBeInTheDocument();
    expect(screen.getByText("5%")).toBeInTheDocument();

    // needs_rebalance badge 應顯示 "Needs Rebalance" (badge + card description label)
    const needsRebalanceElements = screen.getAllByText("Needs Rebalance");
    expect(needsRebalanceElements.length).toBeGreaterThanOrEqual(1);
    // 確認 Badge 存在
    const badge = needsRebalanceElements.find(
      (el) => el.getAttribute("data-slot") === "badge"
    );
    expect(badge).toBeDefined();

    // 配置偏差表格與建議卡片中應顯示資產類別（可能重複出現）
    expect(screen.getAllByText("TW Stock").length).toBeGreaterThanOrEqual(1);
    expect(screen.getAllByText("US Stock").length).toBeGreaterThanOrEqual(1);

    // 建議卡片
    expect(screen.getAllByText("Sell").length).toBeGreaterThanOrEqual(1);
    expect(screen.getAllByText("Buy").length).toBeGreaterThanOrEqual(1);
    expect(
      screen.getByText("台股配置超過目標 10%，建議賣出")
    ).toBeInTheDocument();
  });

  it("should show balanced badge when needs_rebalance is false", async () => {
    server.use(
      http.get(`${API_BASE}/api/rebalance/check`, () => {
        return HttpResponse.json({ data: balancedData });
      })
    );

    renderWithProviders(<RebalancePage />);

    await waitFor(() => {
      expect(screen.getByText("NT$ 1,000,000")).toBeInTheDocument();
    });

    // 應顯示 balanced 狀態（badge + suggestion area 都會出現此文字）
    const balancedElements = screen.getAllByText(
      "Your portfolio is well balanced. No rebalancing is needed at this time."
    );
    expect(balancedElements.length).toBeGreaterThanOrEqual(1);
  });

  it("should show error state on API failure", async () => {
    server.use(
      http.get(`${API_BASE}/api/rebalance/check`, () => {
        return HttpResponse.json(
          { error: { code: "INTERNAL_ERROR", message: "伺服器錯誤" } },
          { status: 500 }
        );
      })
    );

    renderWithProviders(<RebalancePage />);

    await waitFor(() => {
      expect(screen.getByText("Failed to load data")).toBeInTheDocument();
    });

    // 應顯示重新載入按鈕
    expect(screen.getByText("Reload")).toBeInTheDocument();
  });
});
