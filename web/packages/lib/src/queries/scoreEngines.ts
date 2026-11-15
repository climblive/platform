import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type {
  CompetitionID,
  ScoreEngineDescriptor,
  ScoreEngineInstanceID,
} from "../models";
import type { StartScoreEngineArguments } from "../models/rest";

export const getScoreEnginesByCompetitionQuery = (
  competitionId: CompetitionID,
) =>
  createQuery(() => ({
    queryKey: ["score-engines", { competitionId }],
    queryFn: async () =>
      ApiClient.getInstance().getScoreEnginesByCompetition(competitionId),
  }));

export const getScoreEnginesQuery = () =>
  createQuery(() => ({
    queryKey: ["score-engines"],
    queryFn: async () => ApiClient.getInstance().getScoreEngines(),
  }));

export const startScoreEngineMutation = (competitionId: string) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (args: StartScoreEngineArguments) =>
      ApiClient.getInstance().startScoreEngine(competitionId, args),
    onSuccess: (newEngineInstanceId) => {
      const queryKey: QueryKey = ["score-engines", { competitionId }];

      const newEngineDescriptor: ScoreEngineDescriptor = {
        competitionId,
        instanceId: newEngineInstanceId,
      };

      const predicate = ({ instanceId }: ScoreEngineDescriptor) =>
        instanceId === newEngineInstanceId;

      client.setQueryData<ScoreEngineDescriptor[]>(queryKey, (oldEngines) =>
        oldEngines?.some(predicate)
          ? oldEngines
          : [...(oldEngines ?? []), newEngineDescriptor],
      );

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
