import { z } from "@climblive/lib/utils";
import type { Category } from "./generated";

export const categorySchema: z.ZodType<Category> = z.object({
  id: z.number(),
  competitionId: z.number(),
  name: z.string(),
  description: z.string().optional(),
  minAge: z.number().int().nonnegative(),
  maxAge: z.number().int().nonnegative(),
  color: z.string().optional(),
  timeBegin: z.coerce.date(),
  timeEnd: z.coerce.date(),
});
