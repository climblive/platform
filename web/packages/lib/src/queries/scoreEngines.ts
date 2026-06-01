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

export const getScoreEnginesQuery = (contestId: ContestID) =>
  createQuery(() => ({
    queryKey: ["score-engines", { contestId }],
    queryFn: async () => ApiClient.getInstance().getScoreEngines(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 0,
    refetchOnWindowFocus: true,
  }));

export const getRunningScoreEnginesQuery = (
  contestIds: ContestID[] | undefined,
) =>
  createQuery(() => ({
    queryKey: ["score-engines", "all", contestIds ?? []],
    queryFn: async () => {
      const scoreEngines = await Promise.all(
        (contestIds ?? []).map(async (contestId) => ({
          contestId,
          instanceIds: await ApiClient.getInstance().getScoreEngines(contestId),
        })),
      );

      return scoreEngines.flatMap(({ contestId, instanceIds }) =>
        instanceIds.map((instanceId) => ({
          contestId,
          instanceId,
        })),
      );
    },
    enabled: contestIds !== undefined,
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 0,
    refetchOnWindowFocus: true,
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

      client.setQueriesData<RunningScoreEngine[]>(
        {
          queryKey: ["score-engines", "all"],
          exact: false,
        },
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

      client.setQueriesData<RunningScoreEngine[]>(
        {
          queryKey: ["score-engines", "all"],
          exact: false,
        },
        (oldEngines) =>
          oldEngines?.filter(({ instanceId }) => instanceId !== variables),
      );
    },
  }));
};
