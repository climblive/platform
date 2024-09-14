export type Tick = {
  id: number;
  contenderId: number;
  timestamp?: string;
  problemId: number;
  top: boolean;
  attemptsTop: number;
  zone: boolean;
  attemptsZone: number;
};
