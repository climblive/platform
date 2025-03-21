import { createMutation, createQuery, useQueryClient, type QueryKey } from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { CompClass, CompClassTemplate } from "../models";
import { HOUR } from "./constants";

export const getCompClassesQuery = (contestId: number) =>
  createQuery({
    queryKey: ["comp-classes", { contestId }],
    queryFn: async () => ApiClient.getInstance().getCompClasses(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
  });

export const createCompClassMutation = (contestId: number) => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: (template: CompClassTemplate) =>
      ApiClient.getInstance().createCompClass(contestId, template),
    onSuccess: (newCompClass) => {
      let queryKey: QueryKey = ["comp-classes", { contestId }];

      client.setQueryData<CompClass[]>(queryKey, (oldCompClasses) => {
        return [...(oldCompClasses ?? []), newCompClass];
      });

      queryKey = ["comp-class", { id: newCompClass.id }];

      client.setQueryData<CompClass>(queryKey, newCompClass);
    },
  });
};
