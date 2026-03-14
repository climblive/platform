import { z } from "@climblive/lib/utils";
import { ownershipDataSchema } from "./common";
import type { Contest } from "./generated";

export const contestSchema: z.ZodType<Contest> = z.object({
  id: z.number(),
  ownership: ownershipDataSchema,
  archived: z.boolean(),
  location: z.string().optional(),
  country: z.string(),
  seriesId: z.number().optional(),
  name: z.string(),
  description: z.string().optional(),
  qualifyingProblems: z.number(),
  finalists: z.number(),
  info: z.string().optional(),
  gracePeriod: z.number(),
  nameRetentionTime: z.number(),
  timeBegin: z.coerce.date().optional(),
  timeEnd: z.coerce.date().optional(),
  created: z.coerce.date(),
  registeredContenders: z.number(),
});
