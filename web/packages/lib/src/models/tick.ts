import * as z from "zod/v4";
import type { Tick } from "./generated";

export const tickSchema: z.ZodType<Tick> = z.object({
  id: z.number(),
  timestamp: z.coerce.date(),
  problemId: z.number(),
  top: z.boolean(),
  attemptsTop: z.number(),
  zone: z.boolean(),
  attemptsZone: z.number(),
});
