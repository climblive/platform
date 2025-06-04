import * as z from "zod/v4";
import type { Contest } from "./generated";

export const contestSchema: z.ZodType<Contest> = z.object({
  id: z.number(),
  location: z.string().optional(),
  seriesId: z.number().optional(),
  name: z.string(),
  description: z.string().optional(),
  qualifyingProblems: z.number(),
  finalists: z.number(),
  rules: z.string().optional(),
  gracePeriod: z.number(),
  timeBegin: z.coerce.date().optional(),
  timeEnd: z.coerce.date().optional(),
});
