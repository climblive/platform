import { z } from "@climblive/lib/utils";
import type { CompClass } from "./generated";

export const compClassSchema: z.ZodType<CompClass> = z.object({
  id: z.number(),
  competitionId: z.number(),
  name: z.string(),
  description: z.string().optional(),
  color: z.string().optional(),
  timeBegin: z.coerce.date(),
  timeEnd: z.coerce.date(),
});
