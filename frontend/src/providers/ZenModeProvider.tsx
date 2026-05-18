"use client";

import { createContext, useContext, useEffect, useState } from "react";

const STORAGE_KEY = "zen-mode";

type ZenModeContextValue = { zen: boolean; toggleZen: () => void };

const ZenModeContext = createContext<ZenModeContextValue | null>(null);

export function ZenModeProvider({ children }: { children: React.ReactNode }) {
  const [zen, setZen] = useState(false);

  useEffect(() => {
    try {
      if (localStorage.getItem(STORAGE_KEY) === "true") setZen(true);
    } catch {
      // localStorage unavailable — stay default off
    }
  }, []);

  const toggleZen = () =>
    setZen((prev) => {
      const next = !prev;
      try {
        localStorage.setItem(STORAGE_KEY, String(next));
      } catch {
        // ignore persistence failure
      }
      return next;
    });

  return (
    <ZenModeContext.Provider value={{ zen, toggleZen }}>
      {children}
    </ZenModeContext.Provider>
  );
}

export function useZenMode(): ZenModeContextValue {
  const ctx = useContext(ZenModeContext);
  if (!ctx) throw new Error("useZenMode must be used within ZenModeProvider");
  return ctx;
}
