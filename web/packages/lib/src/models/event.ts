import { z } from "@climblive/lib/utils";
import type {
  AscentDeregisteredEvent,
  AscentRegisteredEvent,
  CompetitorPublicInfoUpdatedEvent,
  CompetitorResultUpdatedEvent,
  PointValueUpdatedEvent,
  RaffleWinnerDrawnEvent,
  RulesUpdatedEvent,
} from "./generated";
import { pointValueSchema } from "./pointValue";

export const competitorPublicInfoUpdatedEventSchema: z.ZodType<CompetitorPublicInfoUpdatedEvent> =
  z.object({
    competitorId: z.number(),
    categoryId: z.number(),
    name: z.string(),
    withdrawnFromFinals: z.boolean(),
    disqualified: z.boolean(),
    scrubbedAt: z.coerce.date().optional(),
  });

export const competitorResultUpdatedEventSchema: z.ZodType<CompetitorResultUpdatedEvent> =
  z.object({
    timestamp: z.coerce.date(),
    competitorId: z.number(),
    score: z.string(),
    placement: z.number(),
    rankOrder: z.number(),
    finalist: z.boolean(),
  });

export const ascentRegisteredEventSchema: z.ZodType<AscentRegisteredEvent> =
  z.object({
    tickId: z.number(),
    timestamp: z.coerce.date(),
    competitorId: z.number(),
    revision: z.number(),
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
    competitorId: z.number(),
    problemId: z.number(),
  });

export const pointValueUpdatedEventSchema: z.ZodType<PointValueUpdatedEvent> =
  pointValueSchema;

export const raffleWinnerDrawnEventSchema: z.ZodType<RaffleWinnerDrawnEvent> =
  z.object({
    raffleId: z.number(),
    competitorId: z.number(),
    timestamp: z.coerce.date(),
  });

export const rulesUpdatedEventSchema: z.ZodType<RulesUpdatedEvent> = z.object({
  competitionId: z.number(),
  qualifyingProblems: z.number(),
  finalists: z.number(),
  usePoints: z.boolean(),
  pooledPoints: z.boolean(),
  maxAttempts: z.number().int().min(0).max(999),
  pointDeduction: z
    .number()
    .min(0)
    .max(2 ** 31 - 1),
});
