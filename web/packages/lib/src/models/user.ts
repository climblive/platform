import * as z from "zod";
import type { User } from "./generated";

export const userSchema: z.ZodType<User> = z.object({
    id: z.number(),
    username: z.string(),
    admin: z.boolean(),
    organizers: z.array(z.number())
});