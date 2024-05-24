import { createQuery } from "@tanstack/svelte-query";
import { HOUR } from "./constants";
import { ApiClient } from "../Api";

export const getCompClassesQuery = (contestId: number) =>
  createQuery({
    queryKey: ["compClasses", { contestId }],
    queryFn: async () => ApiClient.getInstance().getCompClasses(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
  });
