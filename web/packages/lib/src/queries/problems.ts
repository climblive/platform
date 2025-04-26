import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Problem, ProblemPatch, ProblemTemplate } from "../models";
import { HOUR } from "./constants";

export const getProblemQuery = (problemId: number) =>
  createQuery({
    queryKey: ["problem", { id: problemId }],
    queryFn: async () => ApiClient.getInstance().getProblem(problemId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
  });

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

export const patchProblemMutation = (problemId: number) => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: (template: ProblemPatch) =>
      ApiClient.getInstance().patchProblem(problemId, template),
    onSuccess: (patchedProblem) => {
      let queryKey: QueryKey = ["problems", { contestId: patchedProblem.contestId }];

      client.setQueryData<Problem[]>(queryKey, (oldProblems) => {
        if (oldProblems === undefined) {
          return undefined;
        }

        return oldProblems.map((problem) => {
          if (problem.id === patchedProblem.id) {
            return patchedProblem;
          }

          return problem;
        });
      });

      queryKey = ["problem", { id: problemId }];

      client.setQueryData<Problem>(queryKey, patchedProblem);
    },
  });
};