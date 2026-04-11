import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { http, HttpResponse } from "msw";
import { server } from "@/test/server";
import { useAnalyticsSummary } from "./useAnalytics";
import type { ReactNode } from "react";

const API_BASE = "http://localhost:8080";

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  });
  return ({ children }: { children: ReactNode }) =>
    QueryClientProvider({ client: queryClient, children });
}

describe("useAnalyticsSummary", () => {
  it("returns loading state initially", () => {
    const { result } = renderHook(() => useAnalyticsSummary(), {
      wrapper: createWrapper(),
    });
    expect(result.current.isLoading).toBe(true);
    expect(result.current.data).toBeUndefined();
  });

  it("returns data on success", async () => {
    server.use(
      http.get(`${API_BASE}/api/analytics/summary`, () => {
        return HttpResponse.json({
          data: {
            total_realized_pl: 15000,
            total_realized_pl_pct: 5.2,
            total_cost_basis: 288000,
            total_sell_amount: 303000,
            total_sell_fee: 400,
            transaction_count: 12,
            currency: "TWD",
            time_range: "month",
            start_date: "2025-10-01",
            end_date: "2025-10-31",
          },
        });
      })
    );

    const { result } = renderHook(() => useAnalyticsSummary(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toBeDefined();
    expect(result.current.data!.total_realized_pl).toBe(15000);
    expect(result.current.data!.transaction_count).toBe(12);
  });

  it("returns error on API failure", async () => {
    server.use(
      http.get(`${API_BASE}/api/analytics/summary`, () => {
        return HttpResponse.json(
          { data: null, error: { code: "INTERNAL_ERROR", message: "fail" } },
          { status: 500 }
        );
      })
    );

    const { result } = renderHook(() => useAnalyticsSummary(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(result.current.error).toBeDefined();
    expect(result.current.error!.code).toBe("HTTP_ERROR");
  });

  it("passes time_range as query param", async () => {
    let capturedUrl = "";
    server.use(
      http.get(`${API_BASE}/api/analytics/summary`, ({ request }) => {
        capturedUrl = request.url;
        return HttpResponse.json({
          data: {
            total_realized_pl: 0,
            total_realized_pl_pct: 0,
            total_cost_basis: 0,
            total_sell_amount: 0,
            total_sell_fee: 0,
            transaction_count: 0,
            currency: "TWD",
            time_range: "year",
            start_date: "2025-01-01",
            end_date: "2025-12-31",
          },
        });
      })
    );

    const { result } = renderHook(() => useAnalyticsSummary("year"), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(capturedUrl).toContain("time_range=year");
  });

  it("defaults to month time range", async () => {
    let capturedUrl = "";
    server.use(
      http.get(`${API_BASE}/api/analytics/summary`, ({ request }) => {
        capturedUrl = request.url;
        return HttpResponse.json({
          data: {
            total_realized_pl: 0,
            total_realized_pl_pct: 0,
            total_cost_basis: 0,
            total_sell_amount: 0,
            total_sell_fee: 0,
            transaction_count: 0,
            currency: "TWD",
            time_range: "month",
            start_date: "2025-10-01",
            end_date: "2025-10-31",
          },
        });
      })
    );

    const { result } = renderHook(() => useAnalyticsSummary(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(capturedUrl).toContain("time_range=month");
  });
});
