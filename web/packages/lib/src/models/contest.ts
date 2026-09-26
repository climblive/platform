import { z } from "@climblive/lib/utils";
import { ownershipDataSchema } from "./common";
import type { Contest } from "./generated";

export const contestSchema: z.ZodType<Contest> = z.object({
  id: z.number(),
  ownership: ownershipDataSchema,
  archivedAt: z.coerce.date().optional(),
  location: z.string().optional(),
  country: z.string(),
  seriesId: z.number().optional(),
  name: z.string(),
  description: z.string().optional(),
  qualifyingProblems: z.number(),
  usePoints: z.boolean(),
  pooledPoints: z.boolean(),
  maxAttempts: z.number().int().min(0).max(999),
  pointDeduction: z
    .number()
    .int()
    .min(0)
    .max(2 ** 31 - 1),
  finalists: z.number(),
  info: z.string().optional(),
  gracePeriod: z.number(),
  nameRetentionTime: z.number(),
  timeBegin: z.coerce.date().optional(),
  timeEnd: z.coerce.date().optional(),
  created: z.coerce.date(),
  registeredContenders: z.number(),
});
