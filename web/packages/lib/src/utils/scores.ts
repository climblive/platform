import type { Problem, Tick } from "../models";

export const calculateProblemScore = (
  problem: Problem,
  tick?: Tick,
): number => {
  let pointValue = 0;

  if (tick) {
    pointValue += problem.pointsTop;

    if (problem.flashBonus && tick.attemptsTop === 1) {
      pointValue += problem.flashBonus;
    }
  }

  return pointValue;
};
