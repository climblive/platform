import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Category, CategoryPatch, CategoryTemplate } from "../models";

export const getCategoryQuery = (categoryId: number) =>
  createQuery(() => ({
    queryKey: ["category", { id: categoryId }],
    queryFn: async () => ApiClient.getInstance().getCategory(categoryId),
  }));

export const getCategoriesQuery = (
  competitionId: string,
  options?: Partial<Parameters<typeof createQuery<Category[]>>[0]>,
) => {
  const client = useQueryClient();

  return createQuery(() => ({
    ...options,
    queryKey: ["categories", { competitionId }],
    queryFn: async () => {
      const categories =
        await ApiClient.getInstance().getCategories(competitionId);

      for (const category of categories) {
        client.setQueryData<Category>(
          ["category", { id: category.id }],
          category,
        );
      }

      return categories;
    },
  }));
};

export const createCategoryMutation = (competitionId: string) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (template: CategoryTemplate) =>
      ApiClient.getInstance().createCategory(competitionId, template),
    onSuccess: (newCategory) => {
      let queryKey: QueryKey = ["categories", { competitionId }];

      client.setQueryData<Category[]>(queryKey, (oldCategories) => {
        return [...(oldCategories ?? []), newCategory];
      });

      queryKey = ["category", { id: newCategory.id }];

      client.setQueryData<Category>(queryKey, newCategory);
    },
  }));
};

export const patchCategoryMutation = (categoryId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (patch: CategoryPatch) =>
      ApiClient.getInstance().patchCategory(categoryId, patch),
    onSuccess: (patchedCategory) => {
      let queryKey: QueryKey = [
        "categories",
        { competitionId: patchedCategory.competitionId },
      ];

      client.setQueryData<Category[]>(queryKey, (oldCategories) => {
        if (oldCategories === undefined) {
          return undefined;
        }

        return oldCategories.map((category) => {
          if (category.id === patchedCategory.id) {
            return patchedCategory;
          }

          return category;
        });
      });

      queryKey = ["category", { id: categoryId }];

      client.setQueryData<Category>(queryKey, patchedCategory);
    },
  }));
};

export const deleteCategoryMutation = (categoryId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: () => ApiClient.getInstance().deleteCategory(categoryId),
    onSuccess: () => {
      let queryKey: QueryKey = ["categories"];

      client.setQueriesData<Category[]>(
        { queryKey, exact: false },
        (oldCategories) => {
          if (oldCategories === undefined) {
            return undefined;
          }

          return oldCategories.filter(({ id }) => id !== categoryId);
        },
      );

      queryKey = ["category", { id: categoryId }];

      client.removeQueries({ queryKey });
    },
  }));
};
