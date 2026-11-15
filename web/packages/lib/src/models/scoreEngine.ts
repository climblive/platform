import { z } from "@climblive/lib/utils";
import { type ScoreEngineDescriptor } from "./generated";

export const scoreEngineDescriptorSchema: z.ZodType<ScoreEngineDescriptor> =
  z.object({
    competitionId: z.string().uuid(),
    instanceId: z.string().uuid(),
  });
