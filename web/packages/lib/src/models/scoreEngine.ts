import { z } from "@climblive/lib/utils";
import { type ScoreEngineDescriptor } from "./generated";

export const scoreEngineDescriptorSchema: z.ZodType<ScoreEngineDescriptor> =
  z.object({
    contestId: z.number(),
    instanceId: z.string().uuid(),
  });
