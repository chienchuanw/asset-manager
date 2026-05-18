import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";

const toggleZen = vi.fn();
vi.mock("@/providers/ZenModeProvider", () => ({
  useZenMode: () => ({ zen: false, toggleZen }),
}));

import { ZenModeToggle } from "./ZenModeToggle";

describe("ZenModeToggle", () => {
  it("calls toggleZen on click", () => {
    render(<ZenModeToggle />);
    fireEvent.click(screen.getByRole("button", { name: /zen/i }));
    expect(toggleZen).toHaveBeenCalledOnce();
  });
});
