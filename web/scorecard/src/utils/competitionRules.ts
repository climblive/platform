import type { Competition } from "@climblive/lib/models";

type CompetitionRules = Pick<
  Competition,
  | "usePoints"
  | "qualifyingProblems"
  | "pooledPoints"
  | "pointDeduction"
  | "maxAttempts"
>;

export const describeCompetitionRules = (rules: CompetitionRules): string => {
  if (!rules.usePoints) {
    return "";
  }

  const sentences: string[] = [];

  if (rules.qualifyingProblems === 1) {
    sentences.push("Only the hardest problem counts.");
  } else if (rules.qualifyingProblems > 1) {
    sentences.push(
      `Only the hardest ${rules.qualifyingProblems} problems count.`,
    );
  }

  if (rules.pooledPoints) {
    sentences.push(
      "Points are shared among the competitors that complete a problem.",
    );
  }

  if (rules.pointDeduction > 0) {
    sentences.push(
      `There is a ${rules.pointDeduction} point penalty for each failed attempt before reaching the scored top or zone.`,
    );
  }

  if (rules.maxAttempts === 1) {
    sentences.push("Only one attempt is allowed per problem.");
  } else if (rules.maxAttempts > 1) {
    sentences.push(`Up to ${rules.maxAttempts} attempts allowed per problem.`);
  }

  return sentences.join(" ");
};
