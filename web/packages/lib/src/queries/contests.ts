import { createQuery } from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
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
