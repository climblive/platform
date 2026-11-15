import { z } from "@climblive/lib/utils";

export type ScorecardSession = {
  competitorId: number;
  competitionId: string;
  registrationCode: string;
  expiryTime: Date;
};

export const scorecardSessionSchema: z.ZodType<ScorecardSession> = z.object({
  competitorId: z.number(),
  competitionId: z.string().uuid(),
  registrationCode: z.string().length(8),
  expiryTime: z.coerce.date(),
});
