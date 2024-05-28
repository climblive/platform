import type { CompClass } from "./compClass";

export type Scoreboard = {
  contestId: number;
  scores: {
    compClass: CompClass;
    contenders: ScoreboardContender[];
  }[];
};

export type ScoreboardContender = {
  contenderId: number;
  contenderName: string;
  totalScore: number;
  qualifyingScore: number;
};

export type ScoreboardUpdate = {
  compClassId: number;
  contender: ScoreboardContender;
};
