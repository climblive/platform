import { z } from "@climblive/lib/utils";
import type { Raffle, RaffleWinner } from "./generated";

export const raffleSchema: z.ZodType<Raffle> = z.object({
  id: z.number(),
  competitionId: z.string().uuid(),
});

export const raffleWinnerSchema: z.ZodType<RaffleWinner> = z.object({
  id: z.number(),
  raffleId: z.number(),
  competitorId: z.number(),
  competitorName: z.string(),
  competitorScrubbedAt: z.coerce.date().optional(),
  timestamp: z.coerce.date(),
});
