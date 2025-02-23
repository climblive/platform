import { createMutation, createQuery, useQueryClient, type QueryKey } from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Contest, ContestTemplate } from "../models";
import { HOUR } from "./constants";

export const getContestQuery = (contestId: number) =>
  createQuery({
    queryKey: ["contest", { id: contestId }],
    queryFn: async () => ApiClient.getInstance().getContest(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
  });

export const getContestsByOrganizerQuery = (organizerId: number) =>
  createQuery({
    queryKey: ["contests", { organizerId }],
    queryFn: async () =>
      ApiClient.getInstance().getContestsByOrganizer(organizerId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
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

      client.setQueryData<Contest>(queryKey, newContest)
    },
  });
};