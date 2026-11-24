import type { Category } from "../models/generated";

export const formatCategoryEligibility = ({
  eligibilityBasis,
  minAge,
  maxAge,
  minBirthYear,
  maxBirthYear,
}: Pick<
  Category,
  "eligibilityBasis" | "minAge" | "maxAge" | "minBirthYear" | "maxBirthYear"
>): string => {
  if (eligibilityBasis === "birth_year") {
    if (minBirthYear !== null && maxBirthYear !== null) {
      return minBirthYear === maxBirthYear
        ? `Born ${minBirthYear}`
        : `Born ${minBirthYear}–${maxBirthYear}`;
    }
    if (minBirthYear !== null) {
      return `Born ${minBirthYear} or later`;
    }
    if (maxBirthYear !== null) {
      return `Born ${maxBirthYear} or earlier`;
    }
    return "Any birth year";
  }

  if (minAge !== null && maxAge !== null) {
    return minAge === maxAge ? `Age ${minAge}` : `Ages ${minAge}–${maxAge}`;
  }
  if (minAge !== null) {
    return `Ages ${minAge} or older`;
  }
  if (maxAge !== null) {
    return `Ages ${maxAge} or younger`;
  }
  return "All ages";
};

export const ordinalSuperscript = (value: number): string => {
  switch (true) {
    case value === 1:
      return "st";
    case value === 2:
      return "nd";
    case value === 3:
      return "rd";
    case value <= 20:
      return `th`;
    case value % 10 === 1:
      return `st`;
    case value % 10 === 2:
      return `nd`;
    case value % 10 === 3:
      return `rd`;
    default:
      return `th`;
  }
};
