import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";

vi.mock("next-intl", () => ({
  useTranslations: () => (k: string) => k,
}));

import { EmptyState } from "./EmptyState";
import { LoadingState } from "./LoadingState";
import { ErrorState } from "./ErrorState";

describe("state components", () => {
  it("EmptyState renders title", () => {
    render(<EmptyState />);
    expect(screen.getByText("emptyTitle")).toBeInTheDocument();
  });

  it("LoadingState renders the chart variant skeleton", () => {
    const { container } = render(<LoadingState variant="chart" />);
    expect(container.querySelector('[data-slot="skeleton"]')).toBeTruthy();
  });

  it("ErrorState retry button invokes onRetry", () => {
    const onRetry = vi.fn();
    render(<ErrorState onRetry={onRetry} />);
    fireEvent.click(screen.getByRole("button", { name: "retry" }));
    expect(onRetry).toHaveBeenCalledOnce();
  });
});
