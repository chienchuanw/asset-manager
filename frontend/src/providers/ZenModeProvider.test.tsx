import { describe, it, expect, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { ZenModeProvider, useZenMode } from "./ZenModeProvider";

function Probe() {
  const { zen, toggleZen } = useZenMode();
  return <button onClick={toggleZen}>{zen ? "on" : "off"}</button>;
}

describe("ZenModeProvider", () => {
  beforeEach(() => localStorage.clear());

  it("defaults to off when storage empty", () => {
    render(
      <ZenModeProvider>
        <Probe />
      </ZenModeProvider>
    );
    expect(screen.getByRole("button").textContent).toBe("off");
  });

  it("toggles and persists to localStorage", () => {
    render(
      <ZenModeProvider>
        <Probe />
      </ZenModeProvider>
    );
    fireEvent.click(screen.getByRole("button"));
    expect(screen.getByRole("button").textContent).toBe("on");
    expect(localStorage.getItem("zen-mode")).toBe("true");
  });

  it("hydrates from existing localStorage value", () => {
    localStorage.setItem("zen-mode", "true");
    render(
      <ZenModeProvider>
        <Probe />
      </ZenModeProvider>
    );
    expect(screen.getByRole("button").textContent).toBe("on");
  });
});
