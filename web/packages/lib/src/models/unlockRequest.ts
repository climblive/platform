import { z } from "zod";
import {
  type UnlockRequest,
  type UnlockRequestReview,
  type UnlockRequestTemplate,
  UnlockRequestStatusApproved,
  UnlockRequestStatusPending,
  UnlockRequestStatusRejected,
} from "./generated";

export const unlockRequestStatusSchema = z.enum([
  UnlockRequestStatusPending,
  UnlockRequestStatusApproved,
  UnlockRequestStatusRejected,
]);

export const unlockRequestSchema = z.object({
  id: z.number(),
  contestId: z.number(),
  organizerId: z.number(),
  status: unlockRequestStatusSchema,
  createdAt: z.coerce.date(),
  reviewedAt: z.coerce.date().optional(),
}) satisfies z.ZodType<UnlockRequest>;

export const unlockRequestTemplateSchema = z.object({
  contestId: z.number(),
}) satisfies z.ZodType<UnlockRequestTemplate>;

export const unlockRequestReviewSchema = z.object({
  status: unlockRequestStatusSchema,
}) satisfies z.ZodType<UnlockRequestReview>;
