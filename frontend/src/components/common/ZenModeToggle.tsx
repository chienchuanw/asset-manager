"use client";

import { useZenMode } from "@/providers/ZenModeProvider";
import { Eye, EyeOff } from "lucide-react";
import { Button } from "@/components/ui/button";

export function ZenModeToggle() {
  const { zen, toggleZen } = useZenMode();
  return (
    <Button
      variant="ghost"
      size="icon"
      aria-label="Toggle zen mode"
      onClick={toggleZen}
    >
      {zen ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
    </Button>
  );
}
