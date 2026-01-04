import {
  createMutation,
  createQuery,
  useQueryClient,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import {
  type UnlockRequestReview,
  type UnlockRequestTemplate,
} from "../models";

export const getUnlockRequestsByContestQuery = (contestId: number) =>
  createQuery(() => ({
    queryKey: ["unlockRequests", { contestId }],
    queryFn: async () =>
      ApiClient.getInstance().getUnlockRequestsByContest(contestId),
    retry: false,
  }));

export const getUnlockRequestsByOrganizerQuery = (organizerId: number) =>
  createQuery(() => ({
    queryKey: ["unlockRequests", { organizerId }],
    queryFn: async () =>
      ApiClient.getInstance().getUnlockRequestsByOrganizer(organizerId),
    retry: false,
  }));

export const getPendingUnlockRequestsQuery = () =>
  createQuery(() => ({
    queryKey: ["unlockRequests", "pending"],
    queryFn: async () => ApiClient.getInstance().getPendingUnlockRequests(),
    retry: false,
  }));

export const createUnlockRequestMutation = () => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (template: UnlockRequestTemplate) =>
      ApiClient.getInstance().createUnlockRequest(template),
    onSuccess: (newRequest) => {
      client.invalidateQueries({
        queryKey: ["unlockRequests", { contestId: newRequest.contestId }],
      });
      client.invalidateQueries({
        queryKey: ["unlockRequests", { organizerId: newRequest.organizerId }],
      });
      client.invalidateQueries({
        queryKey: ["contest", { id: newRequest.contestId }],
      });
    },
  }));
};

export const reviewUnlockRequestMutation = (requestId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (review: UnlockRequestReview) =>
      ApiClient.getInstance().reviewUnlockRequest(requestId, review),
    onSuccess: (updatedRequest) => {
      client.invalidateQueries({
        queryKey: ["unlockRequests", { contestId: updatedRequest.contestId }],
      });
      client.invalidateQueries({
        queryKey: [
          "unlockRequests",
          { organizerId: updatedRequest.organizerId },
        ],
      });
      client.invalidateQueries({
        queryKey: ["unlockRequests", "pending"],
      });
      client.invalidateQueries({
        queryKey: ["contest", { id: updatedRequest.contestId }],
      });
    },
  }));
};
