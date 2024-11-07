import * as z from "zod";

export type ContenderPublicInfoUpdatedEvent = {
  contenderId: number;
  compClassId: number;
  publicName?: string;
  clubName?: string;
  withdrawnFromFinals: boolean;
  disqualified: boolean;
};

export const contenderPublicInfoUpdatedEventSchema: z.ZodType<ContenderPublicInfoUpdatedEvent> = z.object({
  contenderId: z.number(),
  compClassId: z.number(),
  publicName: z.string().optional(),
  clubName: z.string().optional(),
  withdrawnFromFinals: z.boolean(),
  disqualified: z.boolean(),
});

export type ContenderScoreUpdatedEvent = {
  timestamp: Date;
  contenderId: number;
  score: number;
  placement?: number;
  rankOrder: number;
  finalist: boolean;
};

export const contenderScoreUpdatedEventSchema: z.ZodType<ContenderScoreUpdatedEvent> = z.object({
  timestamp: z.coerce.date(),
  contenderId: z.number(),
  score: z.number(),
  placement: z.number().optional(),
  rankOrder: z.number(),
  finalist: z.boolean(),
});