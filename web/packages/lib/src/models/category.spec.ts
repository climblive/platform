import { expect, test } from "vitest";
import { categoryEligibilitySchema } from "./category";

const ageRange = {
  eligibilityBasis: "age",
  minAge: 18,
  maxAge: 99,
  minBirthYear: null,
  maxBirthYear: null,
};

test.each([
  ageRange,
  { ...ageRange, minAge: 0 },
  { ...ageRange, minAge: null },
  { ...ageRange, maxAge: null },
  { ...ageRange, minAge: null, maxAge: null },
  { ...ageRange, maxAge: 18 },
  {
    ...ageRange,
    eligibilityBasis: "birth_year",
    minAge: null,
    maxAge: null,
    minBirthYear: 2010,
    maxBirthYear: 2014,
  },
  {
    ...ageRange,
    eligibilityBasis: "birth_year",
    minAge: null,
    maxAge: null,
    maxBirthYear: 2009,
  },
  {
    ...ageRange,
    eligibilityBasis: "birth_year",
    minAge: null,
    maxAge: null,
    minBirthYear: 2012,
  },
  { ...ageRange, eligibilityBasis: "birth_year", minAge: null, maxAge: null },
])("accepts eligibility %j", (eligibility) => {
  expect(categoryEligibilitySchema.safeParse(eligibility).success).toBe(true);
});

test.each([
  { ...ageRange, eligibilityBasis: "unknown" },
  { ...ageRange, eligibilityBasis: null },
  { ...ageRange, eligibilityBasis: undefined },
  { ...ageRange, minAge: -1 },
  { ...ageRange, maxAge: -1 },
  { ...ageRange, minAge: 1.5 },
  { ...ageRange, maxAge: 2 ** 31 },
  { ...ageRange, minAge: 2 ** 31 },
  { ...ageRange, minAge: 100 },
  { ...ageRange, minBirthYear: 2010 },
  { ...ageRange, maxBirthYear: 2014 },
  { ...ageRange, eligibilityBasis: "birth_year" },
  {
    ...ageRange,
    eligibilityBasis: "birth_year",
    minAge: null,
    maxAge: null,
    minBirthYear: 2014,
    maxBirthYear: 2010,
  },
  {
    ...ageRange,
    eligibilityBasis: "birth_year",
    minAge: null,
    maxAge: null,
    minBirthYear: 0,
  },
  {
    ...ageRange,
    eligibilityBasis: "birth_year",
    minAge: null,
    maxAge: null,
    maxBirthYear: 10000,
  },
])("rejects eligibility %j", (eligibility) => {
  expect(categoryEligibilitySchema.safeParse(eligibility).success).toBe(false);
});
