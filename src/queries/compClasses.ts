import { createQuery } from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import { HOUR, MINUTE } from "../constants";

export const getCompClassesQuery = (contestId: number) =>
  createQuery({
    queryKey: ["compClasses", { contestId }],
    queryFn: async () => ApiClient.getInstance().getCompClasses(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
  });
