import * as z from "zod";
import type { Score } from "./generated";

export const scoreSchema: z.ZodType<Score> = z.object({
  contenderId: z.number(),
  score: z.number(),
  placement: z.number(),
  rankOrder: z.number(),
  finalist: z.boolean(),
  timestamp: z.coerce.date(),
});
