import type { Competition } from "@climblive/lib/models";

type CompetitionRules = Pick<
  Competition,
  | "usePoints"
  | "qualifyingProblems"
  | "pooledPoints"
  | "pointDeduction"
  | "maxAttempts"
>;

export const describeCompetitionRules = (rules: CompetitionRules): string[] => {
  if (!rules.usePoints) {
    return [];
  }

  const descriptions: string[] = [];

  if (rules.qualifyingProblems === 1) {
    descriptions.push("Hardest problem counts");
  } else if (rules.qualifyingProblems > 1) {
    descriptions.push(`Hardest ${rules.qualifyingProblems} problems count`);
  }

  if (rules.pooledPoints) {
    descriptions.push("Points are pooled");
  }

  if (rules.pointDeduction > 0) {
    descriptions.push(
      `${rules.pointDeduction}-point penalty per failed attempt`,
    );
  }

  if (rules.maxAttempts === 1) {
    descriptions.push("1 attempt per problem");
  } else if (rules.maxAttempts > 1) {
    descriptions.push(`Max ${rules.maxAttempts} attempts per problem`);
  }

  return descriptions;
};
