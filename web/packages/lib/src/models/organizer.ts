import { z } from "@climblive/lib/utils";
import type { Organizer } from "./generated";

export const organizerSchema: z.ZodType<Organizer> = z.object({
  id: z.number(),
  name: z.string(),
});
