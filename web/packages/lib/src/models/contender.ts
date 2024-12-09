import * as z from "zod";
import { scoreSchema, type Score } from "./score";

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
  score?: Score;
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
  score: scoreSchema.optional(),
});
