import { createQuery } from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import { HOUR, MINUTE } from "../constants";

export const getProblemsQuery = (contestId: number) =>
  createQuery({
    queryKey: ["problems", { contestId }],
    queryFn: async () => ApiClient.getInstance().getProblems(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
  });
