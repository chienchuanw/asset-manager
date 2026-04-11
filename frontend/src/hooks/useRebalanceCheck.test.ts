import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { http, HttpResponse } from "msw";
import { server } from "@/test/server";
import { useRebalanceCheck } from "./useRebalanceCheck";
import type { RebalanceCheck } from "@/types/rebalance";
import type { ReactNode } from "react";

const API_BASE = "http://localhost:8080";

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  });
  return ({ children }: { children: ReactNode }) => {
    const { createElement } = require("react");
    return createElement(
      QueryClientProvider,
      { client: queryClient },
      children
    );
  };
}

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
    {
      asset_type: "crypto",
      target_percent: 20,
      current_percent: 15,
      deviation: -5,
      deviation_abs: 5,
      exceeds_threshold: true,
      current_value: 150000,
      target_value: 200000,
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
    {
      asset_type: "crypto",
      action: "buy",
      amount: 50000,
      reason: "加密貨幣配置低於目標 5%，建議買入",
    },
  ],
  current_total: 1000000,
};

describe("useRebalanceCheck", () => {
  it("should return loading state initially", () => {
    server.use(
      http.get(`${API_BASE}/api/rebalance/check`, () => {
        // 不回應，讓它保持 loading
        return new Promise(() => {});
      })
    );

    const { result } = renderHook(() => useRebalanceCheck(), {
      wrapper: createWrapper(),
    });

    expect(result.current.isLoading).toBe(true);
    expect(result.current.data).toBeUndefined();
  });

  it("should return data when needs_rebalance is true", async () => {
    server.use(
      http.get(`${API_BASE}/api/rebalance/check`, () => {
        return HttpResponse.json({ data: mockRebalanceData });
      })
    );

    const { result } = renderHook(() => useRebalanceCheck(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toBeDefined();
    expect(result.current.data!.needs_rebalance).toBe(true);
    expect(result.current.data!.deviations).toHaveLength(3);
    expect(result.current.data!.suggestions).toHaveLength(3);
    expect(result.current.data!.current_total).toBe(1000000);
    expect(result.current.data!.threshold).toBe(5);
  });

  it("should return data when needs_rebalance is false", async () => {
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

    server.use(
      http.get(`${API_BASE}/api/rebalance/check`, () => {
        return HttpResponse.json({ data: balancedData });
      })
    );

    const { result } = renderHook(() => useRebalanceCheck(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data!.needs_rebalance).toBe(false);
    expect(result.current.data!.suggestions).toHaveLength(0);
  });

  it("should return error state on API failure", async () => {
    server.use(
      http.get(`${API_BASE}/api/rebalance/check`, () => {
        return HttpResponse.json(
          { error: { code: "INTERNAL_ERROR", message: "伺服器錯誤" } },
          { status: 500 }
        );
      })
    );

    const { result } = renderHook(() => useRebalanceCheck(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isError).toBe(true));

    expect(result.current.error).toBeDefined();
  });
});
