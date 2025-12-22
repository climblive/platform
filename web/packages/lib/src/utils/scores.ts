import type { Problem, Tick } from "../models";

export const calculateProblemScore = (
  problem: Problem,
  tick?: Tick,
): number => {
  let pointValue = 0;

  if (tick?.top) {
    pointValue += problem.pointsTop;

    if (problem.flashBonus && tick.attemptsTop === 1) {
      pointValue += problem.flashBonus;
    }
  } else if (tick?.zone2 && problem.pointsZone2) {
    pointValue += problem.pointsZone2;
  } else if (tick?.zone1 && problem.pointsZone1) {
    pointValue += problem.pointsZone1;
  }

  return pointValue;
};
