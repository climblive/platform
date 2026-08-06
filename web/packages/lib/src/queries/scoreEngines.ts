import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type {
  ContestID,
  ScoreEngineDescriptor,
  ScoreEngineInstanceID,
} from "../models";
import type { StartScoreEngineArguments } from "../models/rest";
import { HOUR } from "./constants";

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
    onSuccess: (newEngineInstanceId) => {
      const queryKey: QueryKey = ["score-engines", { contestId }];

      const newEngineDescriptor: ScoreEngineDescriptor = {
        contestId,
        instanceId: newEngineInstanceId,
      };

      client.setQueryData<ScoreEngineDescriptor[]>(queryKey, (oldEngines) =>
        oldEngines
          ? [...oldEngines, newEngineDescriptor]
          : [newEngineDescriptor],
      );

      const predicate = ({ instanceId }: ScoreEngineDescriptor) =>
        instanceId === newEngineInstanceId;

      client.setQueryData<ScoreEngineDescriptor[]>(
        ["score-engines"],
        (oldEngines) => {
          if (oldEngines?.some(predicate)) {
            return oldEngines.map((engine) =>
              predicate(engine) ? newEngineDescriptor : engine,
            );
          }

          return [...(oldEngines ?? []), newEngineDescriptor];
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

      client.setQueriesData<ScoreEngineDescriptor[]>(
        {
          queryKey,
          exact: false,
        },
        (oldEngines) =>
          oldEngines
            ? oldEngines.filter(({ instanceId }) => instanceId !== variables)
            : undefined,
      );
    },
  }));
};
