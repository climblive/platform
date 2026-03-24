import { z } from "@climblive/lib/utils";
import type { OrganizerInvite } from "./generated";

export const organizerInviteSchema: z.ZodType<OrganizerInvite> = z.object({
  id: z.string().uuid(),
  organizerId: z.number(),
  organizerName: z.string(),
  expiresAt: z.coerce.date(),
});
