import { z } from "@climblive/lib/utils";
import type { PointValue } from "./generated";

export const pointValueSchema: z.ZodType<PointValue> = z.object({
  contenderId: z.number(),
  problemId: z.number(),
  current: z.number(),
  maximum: z.number(),
});
