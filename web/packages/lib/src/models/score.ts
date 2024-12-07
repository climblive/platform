import * as z from "zod";

export type Score = {
  contenderId: number;
  score: number;
  placement?: number;
  rankOrder: number;
  finalist: boolean;
  scoreUpdated?: Date;
};

export const scoreSchema: z.ZodType<Score> = z.object({
  contenderId: z.number(),
  score: z.number(),
  placement: z.number().optional(),
  rankOrder: z.number(),
  finalist: z.boolean(),
  scoreUpdated: z.coerce.date().optional(),
});
