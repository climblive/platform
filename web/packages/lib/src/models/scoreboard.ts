import { z } from "@climblive/lib/utils";
import type { ScoreboardEntry } from "./generated";
import { scoreSchema } from "./score";

export const scoreboardEntrySchema: z.ZodType<ScoreboardEntry> = z.object({
  contenderId: z.number(),
  compClassId: z.number(),
  name: z.string(),
  withdrawnFromFinals: z.boolean(),
  disqualified: z.boolean(),
  scrubbedAt: z.coerce.date().optional(),
  score: scoreSchema.optional(),
});
