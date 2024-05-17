import { createQuery } from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import { HOUR, MINUTE } from "../constants";

export const getContestQuery = (contestId: number) =>
  createQuery({
    queryKey: ["contest", { id: contestId }],
    queryFn: async () => ApiClient.getInstance().getContest(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
  });
