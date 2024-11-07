import * as z from "zod";

export type Tick = {
  id: number;
  timestamp: Date;
  problemId: number;
  top: boolean;
  attemptsTop: number;
  zone: boolean;
  attemptsZone: number;
};

export const tickSchema: z.ZodType<Tick> = z.object({
  id: z.number(),
  timestamp: z.coerce.date(),
  problemId: z.number(),
  top: z.boolean(),
  attemptsTop: z.number(),
  zone: z.boolean(),
  attemptsZone: z.number(),
});
