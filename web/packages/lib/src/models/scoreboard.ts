import * as z from "zod";
import type { ScoreboardEntry } from "./generated";
import { scoreSchema } from "./score";

export const scoreboardEntrySchema: z.ZodType<ScoreboardEntry> = z.object({
  contenderId: z.number(),
  compClassId: z.number(),
  name: z.string(),
  withdrawnFromFinals: z.boolean(),
  disqualified: z.boolean(),
  score: scoreSchema.optional(),
});
