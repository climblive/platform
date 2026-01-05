import * as z from "zod/v4";
import type { Raffle, RaffleWinner } from "./generated";

export const raffleSchema: z.ZodType<Raffle> = z.object({
  id: z.number(),
  contestId: z.number(),
});

export const raffleWinnerSchema: z.ZodType<RaffleWinner> = z.object({
  id: z.number(),
  raffleId: z.number(),
  contenderId: z.number(),
  contenderName: z.string(),
  timestamp: z.coerce.date(),
});
