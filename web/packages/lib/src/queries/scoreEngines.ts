import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { ContestID, ScoreEngineInstanceID } from "../models";
import type { StartScoreEngineArguments } from "../models/rest";
import { HOUR } from "./constants";

export type RunningScoreEngine = {
  contestId: ContestID;
  instanceId: ScoreEngineInstanceID;
};

export const getScoreEnginesByContestQuery = (contestId: ContestID) =>
  createQuery(() => ({
    queryKey: ["score-engines", { contestId }],
    queryFn: async () =>
      ApiClient.getInstance().getScoreEnginesByContest(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 0,
    refetchOnWindowFocus: true,
  }));

export const getScoreEnginesQuery = () =>
  createQuery(() => ({
    queryKey: ["score-engines"],
    queryFn: async () => ApiClient.getInstance().getScoreEngines(),
  }));

export const startScoreEngineMutation = (contestId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (args: StartScoreEngineArguments) =>
      ApiClient.getInstance().startScoreEngine(contestId, args),
    onSuccess: (newEngine) => {
      const queryKey: QueryKey = ["score-engines", { contestId }];

      client.setQueryData<ScoreEngineInstanceID[]>(queryKey, (oldEngines) =>
        oldEngines ? [...oldEngines, newEngine] : [newEngine],
      );

      client.setQueryData<RunningScoreEngine[]>(
        ["score-engines"],
        (oldEngines) => {
          if (oldEngines?.some(({ instanceId }) => instanceId === newEngine)) {
            return oldEngines;
          }

          return [...(oldEngines ?? []), { contestId, instanceId: newEngine }];
        },
      );
    },
  }));
};

export const stopScoreEngineMutation = () => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (instanceId: ScoreEngineInstanceID) =>
      ApiClient.getInstance().stopScoreEngine(instanceId),
    onSuccess: (...args) => {
      const [, variables] = args;
      const queryKey = ["score-engines"];
      client.setQueriesData<ScoreEngineInstanceID[]>(
        {
          queryKey,
          exact: false,
        },
        (oldEngines) =>
          oldEngines
            ? oldEngines.filter((instanceId) => instanceId !== variables)
            : undefined,
      );

      client.setQueryData<RunningScoreEngine[]>(
        ["score-engines"],
        (oldEngines) =>
          oldEngines?.filter(({ instanceId }) => instanceId !== variables),
      );
    },
  }));
};
