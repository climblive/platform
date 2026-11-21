import { expect, test } from "vitest";
import { formatAgeSpan } from "./text";

test.each([
  [18, 99, "18–99"],
  [0, 15, "0–15"],
  [16, 16, "16–16"],
])("formats ages %i to %i as %s", (minAge, maxAge, expected) => {
  expect(formatAgeSpan(minAge, maxAge)).toBe(expected);
});
