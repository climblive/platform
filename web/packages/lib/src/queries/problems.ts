import {
  createMutation,
  createQuery,
  QueryClient,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Problem, ProblemPatch, ProblemTemplate } from "../models";

export const getProblemQuery = (problemId: number) =>
  createQuery(() => ({
    queryKey: ["problem", { id: problemId }],
    queryFn: async () => ApiClient.getInstance().getProblem(problemId),
  }));

export const getProblemsQuery = (
  competitionId: string,
  options?: Partial<Parameters<typeof createQuery<Problem[]>>[0]>,
) =>
  createQuery(() => ({
    ...options,
    queryKey: ["problems", { competitionId }],
    queryFn: async () =>
      ApiClient.getInstance().getProblemsByCompetition(competitionId),
  }));

export const createProblemMutation = (competitionId: string) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (template: ProblemTemplate) =>
      ApiClient.getInstance().createProblem(competitionId, template),
    onSuccess: (newProblem) => {
      let queryKey: QueryKey = ["problems", { competitionId }];

      client.setQueryData<Problem[]>(queryKey, (oldProblems) => {
        return [...(oldProblems ?? []), newProblem];
      });

      queryKey = ["problem", { id: newProblem.id }];

      client.setQueryData<Problem>(queryKey, newProblem);
    },
  }));
};

export const patchProblemMutation = (problemId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (template: ProblemPatch) =>
      ApiClient.getInstance().patchProblem(problemId, template),
    onSuccess: (patchedProblem) => {
      let queryKey: QueryKey = [
        "problems",
        { competitionId: patchedProblem.competitionId },
      ];

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
  }));
};

export const deleteProblemMutation = (problemId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: () => ApiClient.getInstance().deleteProblem(problemId),
    onSuccess: () => {
      let queryKey: QueryKey = ["problems"];

      client.setQueriesData<Problem[]>(
        { queryKey, exact: false },
        (oldProblems) => {
          if (oldProblems === undefined) {
            return undefined;
          }

          return oldProblems.filter(({ id }) => id !== problemId);
        },
      );

      queryKey = ["problem", { id: problemId }];

      client.removeQueries({ queryKey });
    },
  }));
};

export const refetchProblems = (
  queryClient: QueryClient,
  competitionId: string,
) => {
  queryClient.refetchQueries({ queryKey: ["problems", { competitionId }] });
};
