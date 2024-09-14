export type Problem = {
  id: number;
  contestId: number;
  number: number;
  holdColorPrimary: string;
  holdColorSecondary?: string;
  name?: string;
  description?: string;
  pointsTop: number;
  pointsZone: number;
  flashBonus?: number;
};
