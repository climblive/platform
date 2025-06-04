import * as z from "zod/v4";
import type { ScoreboardEntry } from "./generated";
import { scoreSchema } from "./score";

export const scoreboardEntrySchema: z.ZodType<ScoreboardEntry> = z.object({
  contenderId: z.number(),
  compClassId: z.number(),
  publicName: z.string(),
  clubName: z.string().optional(),
  withdrawnFromFinals: z.boolean(),
  disqualified: z.boolean(),
  score: scoreSchema.optional(),
});
