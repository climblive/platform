export type Contender = {
  id: number;
  contestId: number;
  compClassId?: number;
  registrationCode: string;
  name?: string;
  publicName?: string;
  clubName?: string;
  entered?: string;
  withdrawnFromFinals: boolean;
  disqualified: boolean;
  score: number;
  placement?: number;
  rankOrder: number;
  finalist: boolean;
  scoreUpdated?: string;
};
