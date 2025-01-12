import { createQuery } from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import { HOUR } from "./constants";

export const getCompClassesQuery = (contestId: number) =>
  createQuery({
    queryKey: ["comp-classes", { contestId }],
    queryFn: async () => ApiClient.getInstance().getCompClasses(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
  });
