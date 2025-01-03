import * as z from "zod";
import type { ContenderPublicInfoUpdatedEvent, ContenderScoreUpdatedEvent } from "./generated";

export const contenderPublicInfoUpdatedEventSchema: z.ZodType<ContenderPublicInfoUpdatedEvent> =
  z.object({
    contenderId: z.number(),
    compClassId: z.number(),
    publicName: z.string(),
    clubName: z.string().optional(),
    withdrawnFromFinals: z.boolean(),
    disqualified: z.boolean(),
  });

export const contenderScoreUpdatedEventSchema: z.ZodType<ContenderScoreUpdatedEvent> =
  z.object({
    timestamp: z.coerce.date(),
    contenderId: z.number(),
    score: z.number(),
    placement: z.number(),
    rankOrder: z.number(),
    finalist: z.boolean(),
  });
