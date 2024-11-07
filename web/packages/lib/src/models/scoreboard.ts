import * as z from "zod";

export type ScoreboardEntry = {
  contenderId: number;
  compClassId: number;
  publicName?: string;
  clubName?: string;
  withdrawnFromFinals: boolean;
  disqualified: boolean;
  score: number;
  placement?: number;
  rankOrder: number;
  scoreUpdated?: Date;
  finalist: boolean;
};

export const scoreboardEntrySchema: z.ZodType<ScoreboardEntry> = z.object({
  contenderId: z.number(),
  compClassId: z.number(),
  publicName: z.string().optional(),
  clubName: z.string().optional(),
  withdrawnFromFinals: z.boolean(),
  disqualified: z.boolean(),
  score: z.number(),
  placement: z.number().optional(),
  rankOrder: z.number(),
  scoreUpdated: z.coerce.date().optional(),
  finalist: z.boolean(),
});
