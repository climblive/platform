import * as z from "zod/v4";
import type {
  AscentDeregisteredEvent,
  AscentRegisteredEvent,
  ContenderPublicInfoUpdatedEvent,
  ContenderScoreUpdatedEvent,
  RaffleWinnerDrawnEvent,
} from "./generated";

export const contenderPublicInfoUpdatedEventSchema: z.ZodType<ContenderPublicInfoUpdatedEvent> =
  z.object({
    contenderId: z.number(),
    compClassId: z.number(),
    name: z.string(),
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
    zone1: z.boolean(),
    attemptsZone1: z.number(),
    zone2: z.boolean(),
    attemptsZone2: z.number(),
    top: z.boolean(),
    attemptsTop: z.number(),
  });

export const ascentDeregisteredEventSchema: z.ZodType<AscentDeregisteredEvent> =
  z.object({
    tickId: z.number(),
    contenderId: z.number(),
    problemId: z.number(),
  });

export const raffleWinnerDrawnEventSchema: z.ZodType<RaffleWinnerDrawnEvent> =
  z.object({
    winnerId: z.number(),
    raffleId: z.number(),
    contenderId: z.number(),
    contenderName: z.string(),
    timestamp: z.coerce.date(),
  });
