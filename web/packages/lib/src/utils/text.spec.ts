import { expect, test } from "vitest";
import { formatCategoryEligibility } from "./text";

test.each([
  [18, 99, "Ages 18–99"],
  [0, 15, "Ages 0–15"],
  [16, 16, "Age 16"],
  [18, null, "Ages 18 or older"],
  [null, 15, "Ages 15 or younger"],
  [0, null, "Ages 0 or older"],
  [null, 0, "Ages 0 or younger"],
  [null, null, "All ages"],
])("formats ages %s to %s as %s", (minAge, maxAge, expected) => {
  expect(
    formatCategoryEligibility({
      eligibilityBasis: "age",
      minAge,
      maxAge,
      minBirthYear: null,
      maxBirthYear: null,
    }),
  ).toBe(expected);
});

test.each([
  [2010, 2014, "Born 2010–2014"],
  [2010, 2010, "Born 2010"],
  [2012, null, "Born 2012 or later"],
  [null, 2009, "Born 2009 or earlier"],
  [null, null, "Any birth year"],
])(
  "formats birth years %s to %s as %s",
  (minBirthYear, maxBirthYear, expected) => {
    expect(
      formatCategoryEligibility({
        eligibilityBasis: "birth_year",
        minAge: null,
        maxAge: null,
        minBirthYear,
        maxBirthYear,
      }),
    ).toBe(expected);
  },
);
