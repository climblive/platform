import * as z from "zod";
import type {
  AscentDeregisteredEvent,
  AscentRegisteredEvent,
  ContenderPublicInfoUpdatedEvent,
  ContenderScoreUpdatedEvent,
} from "./generated";

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

export const ascentRegisteredEventSchema: z.ZodType<AscentRegisteredEvent> =
  z.object({
    tickId: z.number(),
    timestamp: z.coerce.date(),
    contenderId: z.number(),
    problemId: z.number(),
    top: z.boolean(),
    attemptsTop: z.number(),
    zone: z.boolean(),
    attemptsZone: z.number(),
  });

export const ascentDeregisteredEventSchema: z.ZodType<AscentDeregisteredEvent> =
  z.object({
    tickId: z.number(),
    contenderId: z.number(),
    problemId: z.number(),
  });