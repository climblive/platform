import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { CompClass, CompClassPatch, CompClassTemplate } from "../models";

export const getCompClassQuery = (compClassId: number) =>
  createQuery(() => ({
    queryKey: ["comp-class", { id: compClassId }],
    queryFn: async () => ApiClient.getInstance().getCompClass(compClassId),
  }));

export const getCompClassesQuery = (
  competitionId: number,
  options?: Partial<Parameters<typeof createQuery<CompClass[]>>[0]>,
) => {
  const client = useQueryClient();

  return createQuery(() => ({
    ...options,
    queryKey: ["comp-classes", { competitionId }],
    queryFn: async () => {
      const compClasses =
        await ApiClient.getInstance().getCompClasses(competitionId);

      for (const compClass of compClasses) {
        client.setQueryData<CompClass>(
          ["comp-class", { id: compClass.id }],
          compClass,
        );
      }

      return compClasses;
    },
  }));
};

export const createCompClassMutation = (competitionId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (template: CompClassTemplate) =>
      ApiClient.getInstance().createCompClass(competitionId, template),
    onSuccess: (newCompClass) => {
      let queryKey: QueryKey = ["comp-classes", { competitionId }];

      client.setQueryData<CompClass[]>(queryKey, (oldCompClasses) => {
        return [...(oldCompClasses ?? []), newCompClass];
      });

      queryKey = ["comp-class", { id: newCompClass.id }];

      client.setQueryData<CompClass>(queryKey, newCompClass);
    },
  }));
};

export const patchCompClassMutation = (compClassId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (patch: CompClassPatch) =>
      ApiClient.getInstance().patchCompClass(compClassId, patch),
    onSuccess: (patchedCompClass) => {
      let queryKey: QueryKey = [
        "comp-classes",
        { competitionId: patchedCompClass.competitionId },
      ];

      client.setQueryData<CompClass[]>(queryKey, (oldCompClasses) => {
        if (oldCompClasses === undefined) {
          return undefined;
        }

        return oldCompClasses.map((compClass) => {
          if (compClass.id === patchedCompClass.id) {
            return patchedCompClass;
          }

          return compClass;
        });
      });

      queryKey = ["comp-class", { id: compClassId }];

      client.setQueryData<CompClass>(queryKey, patchedCompClass);
    },
  }));
};

export const deleteCompClassMutation = (compClassId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: () => ApiClient.getInstance().deleteCompClass(compClassId),
    onSuccess: () => {
      let queryKey: QueryKey = ["comp-classes"];

      client.setQueriesData<CompClass[]>(
        { queryKey, exact: false },
        (oldCompClasses) => {
          if (oldCompClasses === undefined) {
            return undefined;
          }

          return oldCompClasses.filter(({ id }) => id !== compClassId);
        },
      );

      queryKey = ["comp-class", { id: compClassId }];

      client.removeQueries({ queryKey });
    },
  }));
};
