import { z } from "@climblive/lib/utils";
import type { Problem } from "./generated";

export const problemSchema: z.ZodType<Problem> = z.object({
  id: z.number(),
  contestId: z.number(),
  number: z.number(),
  holdColorPrimary: z.string(),
  holdColorSecondary: z.string().optional(),
  name: z.string().optional(),
  description: z.string().optional(),
  zone1Enabled: z.boolean(),
  zone2Enabled: z.boolean(),
  pointsZone1: z.number().optional(),
  pointsZone2: z.number().optional(),
  pointsTop: z.number(),
  flashBonus: z.number().optional(),
});
