export type ScoreboardEntry = {
  contenderId: number;
  compClassId: number;
  publicName: string;
  clubName: string;
  withdrawnFromFinals: boolean;
  disqualified: boolean;
  score: number;
  placement?: number;
  rankOrder: number;
  scoreUpdated?: string;
  finalist: boolean;
};
