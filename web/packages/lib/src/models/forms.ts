import * as z from "zod";

export type RegistrationFormData = {
  name: string;
  club?: string;
  compClassId: number;
};

export const registrationFormSchema: z.ZodType<RegistrationFormData> = z.object(
  {
    name: z.string().min(1),
    club: z.string().optional(),
    compClassId: z.coerce.number(),
  },
);
