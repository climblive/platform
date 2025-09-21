import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import {
  type Contest,
  type ContestPatch,
  type ContestTemplate,
} from "../models";
import { HOUR } from "./constants";

export const getContestQuery = (contestId: number) =>
  createQuery({
    queryKey: ["contest", { id: contestId }],
    queryFn: async () => ApiClient.getInstance().getContest(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
  });

export const getAllContestsQuery = (
  options?: Partial<Parameters<typeof createQuery<Contest[]>>[0]>,
) =>
  createQuery({
    ...options,
    queryKey: ["contest"],
    queryFn: async () => ApiClient.getInstance().getAllContests(),
  });

export const getContestsByOrganizerQuery = (
  organizerId: number,
  options?: Partial<Parameters<typeof createQuery<Contest[]>>[0]>,
) =>
  createQuery({
    ...options,
    queryKey: ["contests", { organizerId }],
    queryFn: async () =>
      ApiClient.getInstance().getContestsByOrganizer(organizerId),
  });

export const createContestMutation = (organizerId: number) => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: (template: ContestTemplate) =>
      ApiClient.getInstance().createContest(organizerId, template),
    onSuccess: (newContest) => {
      let queryKey: QueryKey = ["contests", { organizerId }];

      client.setQueryData<Contest[]>(queryKey, (oldContests) => {
        return [...(oldContests ?? []), newContest];
      });

      queryKey = ["contest", { id: newContest.id }];

      client.setQueryData<Contest>(queryKey, newContest);
    },
  });
};

export const patchContestMutation = (contestId: number) => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: (patch: ContestPatch) =>
      ApiClient.getInstance().patchContest(contestId, patch),
    onSuccess: (patchedContest) => {
      let queryKey: QueryKey = [
        "contests",
        { organizerId: patchedContest.ownership.organizerId },
      ];

      client.setQueryData<Contest[]>(queryKey, (oldContests) => {
        if (oldContests === undefined) {
          return undefined;
        }

        return oldContests.map((contest) => {
          if (contest.id === patchedContest.id) {
            return patchedContest;
          }

          return contest;
        });
      });

      queryKey = ["contest", { id: contestId }];

      client.setQueryData<Contest>(queryKey, patchedContest);
    },
  });
};

export const duplicateContestMutation = (contestId: number) => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: () => ApiClient.getInstance().duplicateContest(contestId),
    onSuccess: (duplicatedContest) => {
      let queryKey: QueryKey = [
        "contests",
        { organizerId: duplicatedContest.ownership.organizerId },
      ];

      client.setQueryData<Contest[]>(queryKey, (oldContests) => {
        return [...(oldContests ?? []), duplicatedContest];
      });

      queryKey = ["contest", { id: duplicatedContest.id }];

      client.setQueryData<Contest>(queryKey, duplicatedContest);
    },
  });
};
