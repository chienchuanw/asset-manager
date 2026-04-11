import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { http, HttpResponse } from "msw";
import { server } from "@/test/server";
import { useCashFlows } from "./useCashFlows";
import type { ReactNode } from "react";

const API_BASE = "http://localhost:8080";

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  });
  return ({ children }: { children: ReactNode }) =>
    QueryClientProvider({ client: queryClient, children });
}

describe("useCashFlows", () => {
  it("returns loading state initially", () => {
    const { result } = renderHook(() => useCashFlows(), {
      wrapper: createWrapper(),
    });
    expect(result.current.isLoading).toBe(true);
    expect(result.current.data).toBeUndefined();
  });

  it("returns data on success", async () => {
    server.use(
      http.get(`${API_BASE}/api/cash-flows`, () => {
        return HttpResponse.json({
          data: [
            {
              id: "cf-1",
              date: "2025-10-25",
              type: "income",
              category_id: "cat-1",
              amount: 50000,
              currency: "TWD",
              description: "十月薪資",
              note: null,
              source_type: null,
              source_id: null,
              target_type: null,
              target_id: null,
              created_at: "2025-10-25T00:00:00Z",
              updated_at: "2025-10-25T00:00:00Z",
            },
          ],
        });
      })
    );

    const { result } = renderHook(() => useCashFlows(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveLength(1);
    expect(result.current.data![0].id).toBe("cf-1");
    expect(result.current.data![0].amount).toBe(50000);
  });

  it("returns error on API failure", async () => {
    server.use(
      http.get(`${API_BASE}/api/cash-flows`, () => {
        return HttpResponse.json(
          { data: null, error: { code: "INTERNAL_ERROR", message: "fail" } },
          { status: 500 }
        );
      })
    );

    const { result } = renderHook(() => useCashFlows(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(result.current.error).toBeDefined();
    expect(result.current.error!.code).toBe("HTTP_ERROR");
  });

  it("passes filters as query params", async () => {
    let capturedUrl = "";
    server.use(
      http.get(`${API_BASE}/api/cash-flows`, ({ request }) => {
        capturedUrl = request.url;
        return HttpResponse.json({ data: [] });
      })
    );

    const { result } = renderHook(
      () => useCashFlows({ type: "expense", limit: 10 }),
      { wrapper: createWrapper() }
    );

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(capturedUrl).toContain("type=expense");
    expect(capturedUrl).toContain("limit=10");
  });
});
