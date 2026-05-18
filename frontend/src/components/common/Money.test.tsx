import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";

const zenState = { zen: false };
vi.mock("@/providers/ZenModeProvider", () => ({
  useZenMode: () => ({ zen: zenState.zen, toggleZen: vi.fn() }),
}));

import { Money } from "./Money";

describe("Money", () => {
  it("renders formatted currency when zen is off", () => {
    zenState.zen = false;
    render(<Money value={1284500} currency="TWD" />);
    expect(screen.getByText("NT$1,284,500")).toBeInTheDocument();
  });

  it("masks the amount when zen is on", () => {
    zenState.zen = true;
    render(<Money value={1284500} currency="TWD" />);
    expect(screen.queryByText("NT$1,284,500")).not.toBeInTheDocument();
    expect(screen.getByText(/•+/)).toBeInTheDocument();
  });
});
