export type Contest = {
  id: number;
  location?: string;
  seriesId?: number;
  protected: boolean;
  name: string;
  description?: string;
  finalEnabled: boolean;
  qualifyingProblems: number;
  finalists: number;
  rules?: string;
  gracePeriod: number;
};
