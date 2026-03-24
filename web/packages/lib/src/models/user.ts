import { z } from "@climblive/lib/utils";
import type { User } from "./generated";
import { organizerSchema } from "./organizer";

export const userSchema: z.ZodType<User> = z.object({
  id: z.number(),
  username: z.string(),
  admin: z.boolean(),
  organizers: z.array(organizerSchema),
});
