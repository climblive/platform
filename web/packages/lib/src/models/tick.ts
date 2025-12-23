import * as z from "zod";
import type { Tick } from "./generated";

export const tickSchema: z.ZodType<Tick> = z.object({
  id: z.number(),
  timestamp: z.coerce.date(),
  problemId: z.number(),
  zone1: z.boolean(),
  attemptsZone1: z.number(),
  zone2: z.boolean(),
  attemptsZone2: z.number(),
  top: z.boolean(),
  attemptsTop: z.number(),
});
