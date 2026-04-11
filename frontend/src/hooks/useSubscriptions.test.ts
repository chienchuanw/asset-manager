import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { http, HttpResponse } from "msw";
import { server } from "@/test/server";
import { useSubscriptions } from "./useSubscriptions";
import type { ReactNode } from "react";

const API_BASE = "http://localhost:8080";

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  });
  return ({ children }: { children: ReactNode }) =>
    QueryClientProvider({ client: queryClient, children });
}

describe("useSubscriptions", () => {
  it("returns loading state initially", () => {
    const { result } = renderHook(() => useSubscriptions(), {
      wrapper: createWrapper(),
    });
    expect(result.current.isLoading).toBe(true);
    expect(result.current.data).toBeUndefined();
  });

  it("returns data on success", async () => {
    server.use(
      http.get(`${API_BASE}/api/subscriptions`, () => {
        return HttpResponse.json({
          data: [
            {
              id: "sub-1",
              name: "Netflix",
              amount: 390,
              currency: "TWD",
              billing_cycle: "monthly",
              billing_day: 15,
              category_id: "cat-1",
              payment_method: "credit_card",
              account_id: "card-1",
              start_date: "2025-01-15",
              auto_renew: true,
              status: "active",
              created_at: "2025-01-15T00:00:00Z",
              updated_at: "2025-01-15T00:00:00Z",
            },
          ],
        });
      })
    );

    const { result } = renderHook(() => useSubscriptions(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveLength(1);
    expect(result.current.data![0].id).toBe("sub-1");
    expect(result.current.data![0].name).toBe("Netflix");
    expect(result.current.data![0].status).toBe("active");
  });

  it("returns error on API failure", async () => {
    server.use(
      http.get(`${API_BASE}/api/subscriptions`, () => {
        return HttpResponse.json(
          { data: null, error: { code: "INTERNAL_ERROR", message: "fail" } },
          { status: 500 }
        );
      })
    );

    const { result } = renderHook(() => useSubscriptions(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(result.current.error).toBeDefined();
    expect(result.current.error!.code).toBe("HTTP_ERROR");
  });

  it("passes filters as query params", async () => {
    let capturedUrl = "";
    server.use(
      http.get(`${API_BASE}/api/subscriptions`, ({ request }) => {
        capturedUrl = request.url;
        return HttpResponse.json({ data: [] });
      })
    );

    const { result } = renderHook(
      () => useSubscriptions({ status: "active", limit: 10 }),
      { wrapper: createWrapper() }
    );

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(capturedUrl).toContain("status=active");
    expect(capturedUrl).toContain("limit=10");
  });
});
