import { z } from "@climblive/lib/utils";
import type { Category } from "./generated";

export const categorySchema: z.ZodType<Category> = z.object({
  id: z.number(),
  competitionId: z.number(),
  name: z.string(),
  description: z.string().optional(),
  color: z.string().optional(),
  timeBegin: z.coerce.date(),
  timeEnd: z.coerce.date(),
});
