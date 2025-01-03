import * as z from "zod";
import type { Contest } from "./generated";

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
