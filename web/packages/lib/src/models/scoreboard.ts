import { z } from "@climblive/lib/utils";
import type { ScoreboardEntry } from "./generated";
import { resultSchema } from "./result";

export const scoreboardEntrySchema: z.ZodType<ScoreboardEntry> = z.object({
  contenderId: z.number(),
  categoryId: z.number(),
  name: z.string(),
  withdrawnFromFinals: z.boolean(),
  disqualified: z.boolean(),
  scrubbedAt: z.coerce.date().optional(),
  result: resultSchema.optional(),
});
