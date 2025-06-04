import * as z from "zod/v4";
import type { CompClass } from "./generated";

export const compClassSchema: z.ZodType<CompClass> = z.object({
  id: z.number(),
  contestId: z.number(),
  name: z.string(),
  description: z.string().optional(),
  color: z.string().optional(),
  timeBegin: z.coerce.date(),
  timeEnd: z.coerce.date(),
});
