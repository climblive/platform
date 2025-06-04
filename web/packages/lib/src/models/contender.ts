import * as z from "zod/v4";
import type { Contender } from "./generated";
import { scoreSchema } from "./score";

export const contenderSchema: z.ZodType<Contender> = z.object({
  id: z.number(),
  contestId: z.number(),
  compClassId: z.number().optional(),
  registrationCode: z.string(),
  name: z.string().optional(),
  publicName: z.string().optional(),
  clubName: z.string().optional(),
  entered: z.coerce.date().optional(),
  withdrawnFromFinals: z.boolean(),
  disqualified: z.boolean(),
  score: scoreSchema.optional(),
});
