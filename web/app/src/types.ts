import * as z from "zod";

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
  timestamp: z.coerce.date(),
});

