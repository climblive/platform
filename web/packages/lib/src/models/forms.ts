import * as z from "zod";

export type RegistrationFormData = {
  name: string;
  clubName?: string;
  compClassId: number;
  withdrawnFromFinals: boolean;
};

export const registrationFormSchema: z.ZodType<RegistrationFormData> = z.object(
  {
    name: z.string().min(1),
    clubName: z.string().optional(),
    compClassId: z.coerce.number(),
    withdrawnFromFinals: z.coerce.boolean(),
  },
);
