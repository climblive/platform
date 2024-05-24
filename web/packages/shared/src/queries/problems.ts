import { createQuery } from "@tanstack/svelte-query";
import { HOUR } from "./constants";
import { ApiClient } from "../Api";

export const getProblemsQuery = (contestId: number) =>
  createQuery({
    queryKey: ["problems", { contestId }],
    queryFn: async () => ApiClient.getInstance().getProblems(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
  });
