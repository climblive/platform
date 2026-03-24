import { z } from "@climblive/lib/utils";
import type { OwnershipData } from "./generated";

export const ownershipDataSchema: z.ZodType<OwnershipData> = z.object({
  organizerId: z.number(),
});
