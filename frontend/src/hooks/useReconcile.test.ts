import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { http, HttpResponse } from "msw";
import { server } from "@/test/server";
import {
  useReconcileBankAccount,
  useReconcileBatch,
  useReconcileCreditCard,
} from "./useReconcile";
import type { ReactNode } from "react";

const API_BASE = "http://localhost:8080";

function createWrapper(qc?: QueryClient) {
  const queryClient =
    qc ??
    new QueryClient({
      defaultOptions: { queries: { retry: false, gcTime: 0 } },
    });
  return ({ children }: { children: ReactNode }) =>
    QueryClientProvider({ client: queryClient, children });
}

describe("useReconcileBankAccount", () => {
  it("calls reconcile endpoint and returns result", async () => {
    server.use(
      http.post(`${API_BASE}/api/bank-accounts/abc/reconcile`, () =>
        HttpResponse.json({
          data: {
            target_type: "bank_account",
            target_id: "abc",
            delta: 1000,
            cash_flow_id: "cf-1",
            new_balance: 11000,
          },
        })
      )
    );

    const { result } = renderHook(() => useReconcileBankAccount(), {
      wrapper: createWrapper(),
    });
    result.current.mutate({
      id: "abc",
      input: { new_balance: 11000, date: "2026-04-29" },
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.delta).toBe(1000);
    expect(result.current.data?.cash_flow_id).toBe("cf-1");
  });

  it("invalidates bankAccounts / cashFlows / analytics on success", async () => {
    const qc = new QueryClient({
      defaultOptions: { queries: { retry: false, gcTime: 0 } },
    });
    const invalidate = vi.spyOn(qc, "invalidateQueries");

    server.use(
      http.post(`${API_BASE}/api/bank-accounts/abc/reconcile`, () =>
        HttpResponse.json({
          data: {
            target_type: "bank_account",
            target_id: "abc",
            delta: 0,
            new_balance: 1000,
          },
        })
      )
    );

    const { result } = renderHook(() => useReconcileBankAccount(), {
      wrapper: createWrapper(qc),
    });
    result.current.mutate({
      id: "abc",
      input: { new_balance: 1000, date: "2026-04-29" },
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ["bankAccounts"] });
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ["cashFlows"] });
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ["analytics"] });
  });

  it("surfaces server errors", async () => {
    server.use(
      http.post(`${API_BASE}/api/bank-accounts/abc/reconcile`, () =>
        HttpResponse.json(
          { error: { code: "RECONCILE_FAILED", message: "boom" } },
          { status: 500 }
        )
      )
    );
    const { result } = renderHook(() => useReconcileBankAccount(), {
      wrapper: createWrapper(),
    });
    result.current.mutate({
      id: "abc",
      input: { new_balance: 1, date: "2026-04-29" },
    });
    await waitFor(() => expect(result.current.isError).toBe(true));
  });
});

describe("useReconcileCreditCard", () => {
  it("invalidates creditCards / cashFlows / analytics on success", async () => {
    const qc = new QueryClient({
      defaultOptions: { queries: { retry: false, gcTime: 0 } },
    });
    const invalidate = vi.spyOn(qc, "invalidateQueries");

    server.use(
      http.post(`${API_BASE}/api/credit-cards/cc-1/reconcile`, () =>
        HttpResponse.json({
          data: {
            target_type: "credit_card",
            target_id: "cc-1",
            delta: 3000,
            new_used_credit: 8000,
          },
        })
      )
    );

    const { result } = renderHook(() => useReconcileCreditCard(), {
      wrapper: createWrapper(qc),
    });
    result.current.mutate({
      id: "cc-1",
      input: { new_used_credit: 8000, date: "2026-04-29" },
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ["creditCards"] });
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ["cashFlows"] });
  });
});

describe("useReconcileBatch", () => {
  it("invalidates all four query keys", async () => {
    const qc = new QueryClient({
      defaultOptions: { queries: { retry: false, gcTime: 0 } },
    });
    const invalidate = vi.spyOn(qc, "invalidateQueries");

    server.use(
      http.post(`${API_BASE}/api/user-management/reconcile-batch`, () =>
        HttpResponse.json({ data: [] })
      )
    );

    const { result } = renderHook(() => useReconcileBatch(), {
      wrapper: createWrapper(qc),
    });
    result.current.mutate({
      date: "2026-04-29",
      items: [{ target_type: "bank_account", target_id: "x", new_balance: 1 }],
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ["bankAccounts"] });
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ["creditCards"] });
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ["cashFlows"] });
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ["analytics"] });
  });
});
