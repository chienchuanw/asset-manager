import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { http, HttpResponse } from "msw";
import { server } from "@/test/server";
import { useBankAccounts } from "./useBankAccounts";
import type { ReactNode } from "react";

const API_BASE = "http://localhost:8080";

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  });
  return ({ children }: { children: ReactNode }) =>
    QueryClientProvider({ client: queryClient, children });
}

describe("useBankAccounts", () => {
  it("returns loading state initially", () => {
    const { result } = renderHook(() => useBankAccounts(), {
      wrapper: createWrapper(),
    });
    expect(result.current.isLoading).toBe(true);
    expect(result.current.data).toBeUndefined();
  });

  it("returns data on success", async () => {
    server.use(
      http.get(`${API_BASE}/api/bank-accounts`, () => {
        return HttpResponse.json({
          data: [
            {
              id: "ba-1",
              bank_name: "台北富邦銀行",
              account_type: "savings",
              account_number_last4: "1234",
              currency: "TWD",
              balance: 150000,
              note: null,
              created_at: "2025-01-01T00:00:00Z",
              updated_at: "2025-01-01T00:00:00Z",
            },
          ],
        });
      })
    );

    const { result } = renderHook(() => useBankAccounts(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveLength(1);
    expect(result.current.data![0].id).toBe("ba-1");
    expect(result.current.data![0].bank_name).toBe("台北富邦銀行");
    expect(result.current.data![0].balance).toBe(150000);
  });

  it("returns error on API failure", async () => {
    server.use(
      http.get(`${API_BASE}/api/bank-accounts`, () => {
        return HttpResponse.json(
          { data: null, error: { code: "INTERNAL_ERROR", message: "fail" } },
          { status: 500 }
        );
      })
    );

    const { result } = renderHook(() => useBankAccounts(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(result.current.error).toBeDefined();
    expect(result.current.error!.code).toBe("HTTP_ERROR");
  });

  it("passes filters as query params", async () => {
    let capturedUrl = "";
    server.use(
      http.get(`${API_BASE}/api/bank-accounts`, ({ request }) => {
        capturedUrl = request.url;
        return HttpResponse.json({ data: [] });
      })
    );

    const { result } = renderHook(
      () => useBankAccounts({ currency: "USD" }),
      { wrapper: createWrapper() }
    );

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(capturedUrl).toContain("currency=USD");
  });
});
