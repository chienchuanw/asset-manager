import { describe, it, expect } from "vitest";
import { readdirSync, readFileSync, statSync } from "node:fs";
import { join } from "node:path";

// Dark mode relies on semantic design tokens. Hardcoded light Tailwind
// utilities (bg-white, bg-gray-N, text-gray-N, border-gray-N, text-black,
// arbitrary hex utilities) do not adapt to the dark theme. This guard fails
// if any reappear so the regression cannot creep back in.
// Trailing boundary is a negative lookahead, not \b: the arbitrary-hex
// branch ends in "]", a non-word char, so \b after it could never match
// (it would require a following word char, but real code has whitespace or
// a quote there). (?![\w-]) lets every branch terminate correctly.
const BANNED =
  /\b(?:bg-white|text-black|(?:bg|text|border)-gray-\d{2,3}|(?:bg|text|border)-\[#[0-9a-fA-F]{3,8}\])(?![\w-])/;

// Documented, deliberate exceptions (see issue #25 / PR #119).
const EXCEPT = new Set<string>([
  // Arbitrary selectors that neutralize Recharts' own #fff/#ccc defaults.
  "src/components/ui/chart.tsx",
]);

const ROOT = join(__dirname, "..");

function walk(dir: string, acc: string[] = []): string[] {
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    if (statSync(full).isDirectory()) {
      walk(full, acc);
    } else if (/\.tsx?$/.test(entry) && !/\.test\.tsx?$/.test(entry)) {
      acc.push(full);
    }
  }
  return acc;
}

describe("no hardcoded light colors (dark-mode guard)", () => {
  it("every source file uses semantic tokens, not hardcoded light classes", () => {
    const violations: string[] = [];
    for (const file of walk(ROOT)) {
      const rel = file.slice(file.indexOf("src/"));
      if (EXCEPT.has(rel)) continue;
      readFileSync(file, "utf8")
        .split("\n")
        .forEach((line, i) => {
          if (BANNED.test(line)) {
            violations.push(`${rel}:${i + 1}  ${line.trim().slice(0, 100)}`);
          }
        });
    }
    expect(violations, `Hardcoded light classes found:\n${violations.join("\n")}`).toEqual([]);
  });

  it("BANNED regex actually catches every hardcoded form", () => {
    for (const bad of [
      'className="bg-white"',
      "text-black ",
      "bg-gray-50",
      "bg-gray-100",
      "text-gray-700",
      "border-gray-300",
      "bg-[#fff]",
      'text-[#abcdef]',
      "border-[#ccc]/50",
    ]) {
      expect(BANNED.test(bad), `should flag: ${bad}`).toBe(true);
    }
    for (const ok of [
      "bg-card",
      "bg-muted",
      "text-foreground",
      "text-muted-foreground",
      "border-border",
      "bg-gain/10",
      "text-warning",
      "bg-grayish-thing", // not a gray-NN utility
    ]) {
      expect(BANNED.test(ok), `should not flag: ${ok}`).toBe(false);
    }
  });
});
