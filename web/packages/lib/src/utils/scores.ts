import type { Problem } from "../models/problem";
import type { Tick } from "../models/tick";

export const calculateProblemScore = (
  problem: Problem,
  tick?: Tick,
): number => {
  let pointValue = 0;

  if (tick) {
    pointValue += problem.points;

    if (problem.flashBonus && tick.flash === true) {
      pointValue += problem.flashBonus;
    }
  }

  return pointValue;
};
