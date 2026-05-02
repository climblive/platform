import { z } from "@climblive/lib/utils";
import type { HealthStatus, RunnerStatus } from "./generated";

export const runnerStatusSchema: z.ZodType<RunnerStatus> = z.object({
  healthy: z.boolean(),
  checkedAt: z.coerce.date().optional(),
});

export const healthStatusSchema: z.ZodType<HealthStatus> = z.object({
  scoreEngineManager: runnerStatusSchema,
  scoreKeeper: runnerStatusSchema,
  scrubber: runnerStatusSchema,
});
