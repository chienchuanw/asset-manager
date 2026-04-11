import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { http, HttpResponse } from "msw";
import { server } from "@/test/server";
import { useHoldings } from "./useHoldings";
import type { ReactNode } from "react";

const API_BASE = "http://localhost:8080";

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  });
  return ({ children }: { children: ReactNode }) =>
    QueryClientProvider({ client: queryClient, children });
}

describe("useHoldings", () => {
  it("returns loading state initially", () => {
    const { result } = renderHook(() => useHoldings(), {
      wrapper: createWrapper(),
    });
    expect(result.current.isLoading).toBe(true);
    expect(result.current.data).toBeUndefined();
  });

  it("returns data on success", async () => {
    server.use(
      http.get(`${API_BASE}/api/holdings`, () => {
        return HttpResponse.json({
          data: [
            {
              symbol: "2330",
              name: "台積電",
              asset_type: "tw-stock",
              quantity: 10,
              avg_cost: 550,
              avg_cost_original: 550,
              total_cost: 5500,
              current_price: 600,
              currency: "TWD",
              current_price_twd: 600,
              market_value: 6000,
              unrealized_pl: 500,
              unrealized_pl_pct: 9.09,
            },
          ],
          warnings: [],
        });
      })
    );

    const { result } = renderHook(() => useHoldings(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toBeDefined();
    expect(result.current.data!.data).toHaveLength(1);
    expect(result.current.data!.data[0].symbol).toBe("2330");
    expect(result.current.data!.warnings).toEqual([]);
  });

  it("returns warnings when present", async () => {
    server.use(
      http.get(`${API_BASE}/api/holdings`, () => {
        return HttpResponse.json({
          data: [],
          warnings: [
            {
              code: "PRICE_STALE",
              symbol: "2330",
              message: "Price data is stale",
            },
          ],
        });
      })
    );

    const { result } = renderHook(() => useHoldings(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data!.warnings).toHaveLength(1);
    expect(result.current.data!.warnings[0].code).toBe("PRICE_STALE");
  });

  it("returns error on API failure", async () => {
    server.use(
      http.get(`${API_BASE}/api/holdings`, () => {
        return HttpResponse.json(
          { data: null, error: { code: "INTERNAL_ERROR", message: "fail" } },
          { status: 500 }
        );
      })
    );

    const { result } = renderHook(() => useHoldings(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(result.current.error).toBeDefined();
    expect(result.current.error!.code).toBe("HTTP_ERROR");
  });

  it("passes filters as query params", async () => {
    let capturedUrl = "";
    server.use(
      http.get(`${API_BASE}/api/holdings`, ({ request }) => {
        capturedUrl = request.url;
        return HttpResponse.json({ data: [], warnings: [] });
      })
    );

    const { result } = renderHook(
      () => useHoldings({ asset_type: "tw-stock" }),
      { wrapper: createWrapper() }
    );

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(capturedUrl).toContain("asset_type=tw-stock");
  });
});
