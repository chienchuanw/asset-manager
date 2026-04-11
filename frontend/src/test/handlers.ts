import { http, HttpResponse } from "msw";

const API_BASE = "http://localhost:8080";

export const handlers = [
  http.get(`${API_BASE}/api/holdings`, () => {
    return HttpResponse.json({
      data: [],
      warnings: [],
    });
  }),

  http.get(`${API_BASE}/api/transactions`, () => {
    return HttpResponse.json({
      data: [],
      meta: { total: 0, page: 1, page_size: 20, total_pages: 0 },
    });
  }),

  http.get(`${API_BASE}/api/cash-flows`, () => {
    return HttpResponse.json({
      data: [],
      meta: { total: 0, page: 1, page_size: 20, total_pages: 0 },
    });
  }),

  http.get(`${API_BASE}/api/categories`, () => {
    return HttpResponse.json({ data: [] });
  }),

  http.get(`${API_BASE}/api/rebalance/check`, () => {
    return HttpResponse.json({
      data: {
        needs_rebalance: false,
        threshold: 5,
        deviations: [],
        suggestions: [],
        current_total: 0,
      },
    });
  }),
];
