import * as z from "zod";

export type Contender = {
  id: number;
  contestId: number;
  compClassId?: number;
  registrationCode: string;
  name?: string;
  publicName?: string;
  clubName?: string;
  entered?: Date;
  withdrawnFromFinals: boolean;
  disqualified: boolean;
  score: number;
  placement?: number;
  rankOrder: number;
  finalist: boolean;
  scoreUpdated?: Date;
};

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
  score: z.number(),
  placement: z.number().optional(),
  rankOrder: z.number(),
  finalist: z.boolean(),
  scoreUpdated: z.coerce.date().optional(),
});
