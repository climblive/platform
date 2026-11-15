import { z } from "@climblive/lib/utils";
import type { Contender } from "./generated";
import { resultSchema } from "./result";

export const contenderSchema: z.ZodType<Contender> = z.object({
  id: z.number(),
  competitionId: z.number(),
  categoryId: z.number().optional(),
  registrationCode: z.string(),
  name: z.string().optional(),
  entered: z.coerce.date().optional(),
  withdrawnFromFinals: z.boolean(),
  disqualified: z.boolean(),
  scrubbedAt: z.coerce.date().optional(),
  scrubBefore: z.coerce.date().optional(),
  result: resultSchema.optional(),
});
