import * as z from "zod";
import type { Organizer } from "./generated";

export const organizerSchema: z.ZodType<Organizer> = z.object({
    id: z.number(),
    name: z.string(),
});