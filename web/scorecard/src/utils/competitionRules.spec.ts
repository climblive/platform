import { describe, expect, it } from "vitest";
import { describeCompetitionRules } from "./competitionRules";

const defaultRules = {
  usePoints: true,
  qualifyingProblems: 0,
  pooledPoints: false,
  pointDeduction: 0,
  maxAttempts: 0,
};

describe(describeCompetitionRules.name, () => {
  it("describes all enabled rules in order", () => {
    expect(
      describeCompetitionRules({
        usePoints: true,
        qualifyingProblems: 8,
        pooledPoints: true,
        pointDeduction: 0.1,
        maxAttempts: 10,
      }),
    ).toBe(
      "Only the hardest 8 problems count. " +
        "Points are shared among the competitors that complete a problem. " +
        "There is a 0.1 point penalty for each failed attempt before reaching the scored top or zone. " +
        "Up to 10 attempts allowed per problem.",
    );
  });

  it("describes unrestricted rules without zero limits or penalties", () => {
    expect(describeCompetitionRules(defaultRules)).toBe(
      "All problems count. " +
        "Points are not shared between competitors. " +
        "There is no point penalty for failed attempts. " +
        "There is no attempt limit per problem.",
    );
  });

  it.each([
    {
      rule: { qualifyingProblems: 1 },
      sentence: "Only the hardest problem counts.",
    },
    {
      rule: { qualifyingProblems: 8 },
      sentence: "Only the hardest 8 problems count.",
    },
    {
      rule: { pooledPoints: true },
      sentence:
        "Points are shared among the competitors that complete a problem.",
    },
    {
      rule: { pointDeduction: 2.5 },
      sentence:
        "There is a 2.5 point penalty for each failed attempt before reaching the scored top or zone.",
    },
    {
      rule: { maxAttempts: 1 },
      sentence: "Only one attempt is allowed per problem.",
    },
    {
      rule: { maxAttempts: 10 },
      sentence: "Up to 10 attempts allowed per problem.",
    },
  ])("describes an individual rule: $sentence", ({ rule, sentence }) => {
    expect(describeCompetitionRules({ ...defaultRules, ...rule })).toContain(
      sentence,
    );
  });

  it("does not describe point-based options when ranking by attempts", () => {
    expect(
      describeCompetitionRules({
        usePoints: false,
        qualifyingProblems: 8,
        pooledPoints: true,
        pointDeduction: 0.1,
        maxAttempts: 10,
      }),
    ).toBe("All problems count.");
  });
});
