import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { ContestID, ScoreEngineInstanceID } from "../models";
import { HOUR } from "./constants";

export const getScoreEnginesQuery = (contestId: ContestID) =>
  createQuery({
    queryKey: ["score-engines", { contestId }],
    queryFn: async () => ApiClient.getInstance().getScoreEngines(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 0,
    refetchOnWindowFocus: true,
  });

export const startScoreEngineMutation = (contestId: number) => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: () =>
      ApiClient.getInstance().startScoreEngine(contestId),
    onSuccess: (newEngine) => {
      const queryKey: QueryKey = ["score-engines", { contestId }];

      client.setQueryData<ScoreEngineInstanceID[]>(queryKey, (oldEngines) =>
        oldEngines ? [...oldEngines, newEngine] : [newEngine],
      );
    },
  });
};

export const stopScoreEngineMutation = () => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: (instanceId: ScoreEngineInstanceID) => ApiClient.getInstance().stopScoreEngine(instanceId),
    onSuccess: (...args) => {
      const [, variables] = args;
      const queryKey = ["score-engines"];
      client.setQueriesData<ScoreEngineInstanceID[]>(
        {
          queryKey,
          exact: false,
        },
        (oldEngines) =>
          oldEngines ? oldEngines.filter((instanceId) => instanceId !== variables) : undefined,
      );
    },
  });
};
