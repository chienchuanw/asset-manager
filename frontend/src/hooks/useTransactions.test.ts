import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { http, HttpResponse } from "msw";
import { server } from "@/test/server";
import { useTransactions } from "./useTransactions";
import type { ReactNode } from "react";

const API_BASE = "http://localhost:8080";

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  });
  return ({ children }: { children: ReactNode }) =>
    QueryClientProvider({ client: queryClient, children });
}

describe("useTransactions", () => {
  it("returns loading state initially", () => {
    const { result } = renderHook(() => useTransactions(), {
      wrapper: createWrapper(),
    });
    expect(result.current.isLoading).toBe(true);
    expect(result.current.data).toBeUndefined();
  });

  it("returns data on success", async () => {
    server.use(
      http.get(`${API_BASE}/api/transactions`, () => {
        return HttpResponse.json({
          data: [
            {
              id: "txn-1",
              date: "2025-10-23T00:00:00Z",
              asset_type: "tw-stock",
              symbol: "2330",
              name: "台積電",
              type: "buy",
              quantity: 10,
              price: 620,
              amount: 6200,
              fee: 20,
              tax: null,
              currency: "TWD",
              exchange_rate_id: null,
              note: null,
              created_at: "2025-10-23T00:00:00Z",
              updated_at: "2025-10-23T00:00:00Z",
            },
          ],
        });
      })
    );

    const { result } = renderHook(() => useTransactions(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveLength(1);
    expect(result.current.data![0].id).toBe("txn-1");
    expect(result.current.data![0].symbol).toBe("2330");
  });

  it("returns error on API failure", async () => {
    server.use(
      http.get(`${API_BASE}/api/transactions`, () => {
        return HttpResponse.json(
          { data: null, error: { code: "INTERNAL_ERROR", message: "fail" } },
          { status: 500 }
        );
      })
    );

    const { result } = renderHook(() => useTransactions(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(result.current.error).toBeDefined();
    expect(result.current.error!.code).toBe("HTTP_ERROR");
  });

  it("passes filters as query params", async () => {
    let capturedUrl = "";
    server.use(
      http.get(`${API_BASE}/api/transactions`, ({ request }) => {
        capturedUrl = request.url;
        return HttpResponse.json({ data: [] });
      })
    );

    const { result } = renderHook(
      () => useTransactions({ asset_type: "us-stock", limit: 5 }),
      { wrapper: createWrapper() }
    );

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(capturedUrl).toContain("asset_type=us-stock");
    expect(capturedUrl).toContain("limit=5");
  });
});
