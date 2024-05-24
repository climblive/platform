import * as z from "zod";
import type { ScoreboardContender } from "@climblive/shared/models/scoreboard";

export type ScorecardSession = {
  contenderId: number;
  contestId: number;
  registrationCode: string;
  timestamp: Date;
};

export const scorecardSessionSchema: z.ZodType<ScorecardSession> = z.object({
  contenderId: z.number(),
  contestId: z.number(),
  registrationCode: z.string().length(8),
  timestamp: z.coerce.date()
});

export type RankedContender = ScoreboardContender & {
  order: number;
  placement: number;
  finalist: boolean;
};
