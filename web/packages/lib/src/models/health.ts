import { z } from "@climblive/lib/utils";
import type { ServiceStatus } from "./generated";

export const serviceStatusSchema: z.ZodType<ServiceStatus> = z.object({
  name: z.string(),
  healthy: z.boolean(),
  checkedAt: z.coerce.date(),
});
