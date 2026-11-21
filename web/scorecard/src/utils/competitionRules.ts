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
    return "All problems count.";
  }

  const sentences = [
    rules.qualifyingProblems > 0
      ? rules.qualifyingProblems === 1
        ? "Only the hardest problem counts."
        : `Only the hardest ${rules.qualifyingProblems} problems count.`
      : "All problems count.",
    rules.pooledPoints
      ? "Points are shared among the competitors that complete a problem."
      : "Points are not shared between competitors.",
    rules.pointDeduction > 0
      ? `There is a ${rules.pointDeduction} point penalty for each failed attempt before reaching the scored top or zone.`
      : "There is no point penalty for failed attempts.",
    rules.maxAttempts > 0
      ? rules.maxAttempts === 1
        ? "Only one attempt is allowed per problem."
        : `Up to ${rules.maxAttempts} attempts allowed per problem.`
      : "There is no attempt limit per problem.",
  ];

  return sentences.join(" ");
};
