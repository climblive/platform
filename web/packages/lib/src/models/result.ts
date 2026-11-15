import { z } from "@climblive/lib/utils";
import type { Result } from "./generated";

export const resultSchema: z.ZodType<Result> = z.object({
  competitorId: z.number(),
  score: z.string(),
  placement: z.number(),
  rankOrder: z.number(),
  finalist: z.boolean(),
  timestamp: z.coerce.date(),
});
