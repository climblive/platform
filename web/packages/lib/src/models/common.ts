import * as z from "zod";
import type { OwnershipData } from "./generated";

export const ownershipDataSchema: z.ZodType<OwnershipData> = z.object({
  organizerId: z.number(),
});
