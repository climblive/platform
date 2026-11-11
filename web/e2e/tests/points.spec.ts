import { expect, test } from "@playwright/test";
import {
  localizedPointsValue,
  roundPoints,
} from "../../packages/lib/src/utils/points";

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
