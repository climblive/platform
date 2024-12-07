import * as z from "zod";

export type Contest = {
  id: number;
  location?: string;
  seriesId?: number;
  protected: boolean;
  name: string;
  description?: string;
  finalsEnabled: boolean;
  qualifyingProblems: number;
  finalists: number;
  rules?: string;
  gracePeriod: number;
  timeBegin?: Date;
  timeEnd?: Date;
};

export const contestSchema: z.ZodType<Contest> = z.object({
  id: z.number(),
  location: z.string().optional(),
  seriesId: z.number().optional(),
  protected: z.boolean(),
  name: z.string(),
  description: z.string().optional(),
  finalsEnabled: z.boolean(),
  qualifyingProblems: z.number(),
  finalists: z.number(),
  rules: z.string().optional(),
  gracePeriod: z.number(),
  timeBegin: z.coerce.date().optional(),
  timeEnd: z.coerce.date().optional(),
});
