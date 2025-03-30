import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Problem, ProblemTemplate } from "../models";
import { HOUR } from "./constants";

export const getProblemsQuery = (contestId: number) =>
  createQuery({
    queryKey: ["problems", { contestId }],
    queryFn: async () => ApiClient.getInstance().getProblems(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
  });

export const createProblemMutation = (contestId: number) => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: (template: ProblemTemplate) =>
      ApiClient.getInstance().createProblem(contestId, template),
    onSuccess: (newProblem) => {
      let queryKey: QueryKey = ["problems", { contestId }];

      client.setQueryData<Problem[]>(queryKey, (oldProblems) => {
        return [...(oldProblems ?? []), newProblem];
      });

      queryKey = ["problem", { id: newProblem.id }];

      client.setQueryData<Problem>(queryKey, newProblem);
    },
  });
};
