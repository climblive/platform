import * as z from "zod/v4";

export type ScorecardSession = {
  contenderId: number;
  contestId: number;
  registrationCode: string;
  expiryTime: Date;
};

export const scorecardSessionSchema: z.ZodType<ScorecardSession> = z.object({
  contenderId: z.number(),
  contestId: z.number(),
  registrationCode: z.string().length(8),
  expiryTime: z.coerce.date(),
});
