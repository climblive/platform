import { createQuery } from "@tanstack/svelte-query";
import { HOUR } from "./constants";
import { ApiClient } from "../Api";

export const getContestQuery = (contestId: number) =>
  createQuery({
    queryKey: ["contest", { id: contestId }],
    queryFn: async () => ApiClient.getInstance().getContest(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
  });
