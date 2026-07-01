import { z } from "@climblive/lib/utils";
import type { PointValue } from "./generated";

export const pointValueSchema: z.ZodType<PointValue> = z.object({
  contenderId: z.number(),
  problemId: z.number(),
  current: z.number(),
  zone1: z.number(),
  zone2: z.number(),
  top: z.number(),
  flash: z.number(),
});
