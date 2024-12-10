import * as z from "zod";
import { scoreSchema, type Score } from "./score";

export type ScoreboardEntry = {
  contenderId: number;
  compClassId: number;
  publicName?: string;
  clubName?: string;
  withdrawnFromFinals: boolean;
  disqualified: boolean;
  score?: Score;
};

export const scoreboardEntrySchema: z.ZodType<ScoreboardEntry> = z.object({
  contenderId: z.number(),
  compClassId: z.number(),
  publicName: z.string().optional(),
  clubName: z.string().optional(),
  withdrawnFromFinals: z.boolean(),
  disqualified: z.boolean(),
  score: scoreSchema.optional(),
});
