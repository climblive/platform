import * as z from "zod";

export type CompClass = {
  id: number;
  contestId: number;
  name: string;
  description?: string;
  color?: string;
  timeBegin: Date;
  timeEnd: Date;
};

export const compClassSchema: z.ZodType<CompClass> = z.object({
  id: z.number(),
  contestId: z.number(),
  name: z.string(),
  description: z.string().optional(),
  color: z.string().optional(),
  timeBegin: z.coerce.date(),
  timeEnd: z.coerce.date(),
});
