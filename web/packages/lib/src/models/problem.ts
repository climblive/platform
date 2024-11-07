import * as z from "zod";

export type Problem = {
  id: number;
  contestId: number;
  number: number;
  holdColorPrimary: string;
  holdColorSecondary?: string;
  name?: string;
  description?: string;
  pointsTop: number;
  pointsZone: number;
  flashBonus?: number;
};

export const problemSchema: z.ZodType<Problem> = z.object({
  id: z.number(),
  contestId: z.number(),
  number: z.number(),
  holdColorPrimary: z.string(),
  holdColorSecondary: z.string().optional(),
  name: z.string().optional(),
  description: z.string().optional(),
  pointsTop: z.number(),
  pointsZone: z.number(),
  flashBonus: z.number().optional(),
});
