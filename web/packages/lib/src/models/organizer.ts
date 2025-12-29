import * as z from "zod";
import type { Organizer, OrganizerTemplate } from "./generated";

export const organizerSchema: z.ZodType<Organizer> = z.object({
  id: z.number(),
  name: z.string(),
});

export const organizerTemplateSchema: z.ZodType<OrganizerTemplate> = z.object({
  name: z.string(),
});
