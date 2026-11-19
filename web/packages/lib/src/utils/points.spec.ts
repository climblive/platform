import { afterEach, expect, test, vi } from "vitest";
import {
  localizedPointsValue,
  pointsNumberFormat,
  roundPoints,
} from "./points";

afterEach(() => {
  vi.restoreAllMocks();
});

test("rounds fractional sums without losing a tenth", () => {
  expect(roundPoints(33.3 + 8.3)).toBe(41.6);
  expect(roundPoints(33.3 + 33.3 + 33.3)).toBe(99.9);
  expect(roundPoints(0.1 + 0.1 + 0.6)).toBe(0.8);
  expect(roundPoints(0.1 + 0.6 + 0.1)).toBe(0.8);
});

test("formats fractional sums without losing a tenth", () => {
  const expected =
    (41.6).toLocaleString(undefined, { maximumFractionDigits: 1 }) + "p";

  expect(localizedPointsValue(33.3 + 8.3)).toBe(expected);
});

test.each([
  { points: 0, digits: 0 },
  { points: 100, digits: 0 },
  { points: 100.04, digits: 0 },
  { points: 100.05, digits: 1 },
  { points: 100.5, digits: 1 },
])(
  "fixes normal precision to $digits digits for $points points",
  ({ points, digits }) => {
    const formatter = pointsNumberFormat(points);

    expect(formatter.resolvedOptions().minimumFractionDigits).toBe(digits);
    expect(formatter.resolvedOptions().maximumFractionDigits).toBe(digits);
    expect(formatter.format(12.34)).toBe(
      (12.34).toLocaleString(undefined, {
        minimumFractionDigits: digits,
        maximumFractionDigits: digits,
      }),
    );
  },
);

test("does not inspect formatted parts for normal scores", () => {
  const formatToParts = vi.spyOn(Intl.NumberFormat.prototype, "formatToParts");

  pointsNumberFormat(100);
  pointsNumberFormat(100.5);

  expect(formatToParts).not.toHaveBeenCalled();
});

test("reuses the compact precision formatter", () => {
  const formatToParts = vi.spyOn(Intl.NumberFormat.prototype, "formatToParts");

  pointsNumberFormat(10_000, 9_999);
  pointsNumberFormat(12_345, 9_999);

  expect(formatToParts).toHaveBeenCalledTimes(2);
  expect(formatToParts.mock.contexts[0]).toBe(formatToParts.mock.contexts[1]);
});

test.each([10_000, 12_000, 12_345, 999_949.95, 999_950, 1_000_000])(
  "preserves compact formatting and precision for %s points",
  (points) => {
    const rounded = roundPoints(points);
    const expected = new Intl.NumberFormat(undefined, {
      notation: "compact",
      maximumFractionDigits: 1,
    });
    const digits =
      expected.formatToParts(rounded).find((part) => part.type === "fraction")
        ?.value.length ?? 0;
    const formatter = pointsNumberFormat(points, 9_999);

    expect(formatter.format(rounded)).toBe(expected.format(rounded));
    expect(formatter.resolvedOptions().minimumFractionDigits).toBe(digits);
    expect(formatter.resolvedOptions().maximumFractionDigits).toBe(digits);
    expect(localizedPointsValue(points, 9_999)).toBe(
      expected.format(rounded) + "p",
    );
  },
);
