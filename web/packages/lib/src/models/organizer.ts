import * as z from "zod/v4";
import type { Organizer } from "./generated";

export const organizerSchema: z.ZodType<Organizer> = z.object({
  id: z.number(),
  name: z.string(),
});
