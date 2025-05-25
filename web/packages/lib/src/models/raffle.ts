import * as z from "zod";
import type { Raffle } from "./generated";

export const raffleSchema: z.ZodType<Raffle> = z.object({
  id: z.number(),
  contestId: z.number(),
});
