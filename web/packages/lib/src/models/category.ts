import { z } from "@climblive/lib/utils";
import type { Category } from "./generated";

export const categoryEligibilitySchema = z
  .object({
    eligibilityBasis: z.enum(["age", "birth_year"]),
    minAge: z
      .number()
      .int()
      .min(0)
      .max(2 ** 31 - 1)
      .nullable(),
    maxAge: z
      .number()
      .int()
      .min(0)
      .max(2 ** 31 - 1)
      .nullable(),
    minBirthYear: z.number().int().min(1).max(9999).nullable(),
    maxBirthYear: z.number().int().min(1).max(9999).nullable(),
  })
  .check((ctx) => {
    const data = ctx.value;
    const ageBased = data.eligibilityBasis === "age";
    const minimum = ageBased ? data.minAge : data.minBirthYear;
    const maximum = ageBased ? data.maxAge : data.maxBirthYear;

    if (minimum !== null && maximum !== null && minimum > maximum) {
      ctx.issues.push({
        code: "custom",
        message: ageBased
          ? "Maximum age must be at least the minimum age"
          : "Latest birth year must be at least the earliest birth year",
        input: data,
        path: [ageBased ? "maxAge" : "maxBirthYear"],
      });
    }

    const inactiveFields = ageBased
      ? (["minBirthYear", "maxBirthYear"] as const)
      : (["minAge", "maxAge"] as const);

    for (const field of inactiveFields) {
      if (data[field] !== null) {
        ctx.issues.push({
          code: "custom",
          message: "Bounds must be empty for the inactive eligibility basis",
          input: data,
          path: [field],
        });
      }
    }
  });

export const categorySchema: z.ZodType<Category> =
  categoryEligibilitySchema.safeExtend({
    id: z.number(),
    competitionId: z.number(),
    name: z.string(),
    description: z.string().optional(),
    color: z.string().optional(),
    timeBegin: z.coerce.date(),
    timeEnd: z.coerce.date(),
  });
