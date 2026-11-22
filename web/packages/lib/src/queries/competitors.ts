import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type {
  Competitor,
  CompetitorPatch,
  PointValue,
  Result,
} from "../models";
import type { CreateCompetitorsArguments } from "../models/rest";

export const getCompetitorQuery = (competitorId: number) =>
  createQuery(() => ({
    queryKey: ["competitor", { id: competitorId }],
    queryFn: async () => ApiClient.getInstance().getCompetitor(competitorId),
    refetchOnWindowFocus: "always",
  }));

export const getCompetitorsByCompetitionQuery = (competitionId: number) =>
  createQuery(() => ({
    queryKey: ["competitors", { competitionId }],
    queryFn: async () =>
      ApiClient.getInstance().getCompetitorsByCompetition(competitionId),
  }));

export const getPointValuesByCompetitorQuery = (competitorId: number) =>
  createQuery(() => ({
    queryKey: ["point-values", { competitorId }],
    queryFn: async () =>
      ApiClient.getInstance().getPointValuesByCompetitor(competitorId),
    refetchOnWindowFocus: "always",
  }));

export const patchCompetitorMutation = (competitorId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: async (patch: CompetitorPatch) =>
      ApiClient.getInstance().patchCompetitor(competitorId, patch),
    onSuccess: (updatedCompetitor) => {
      const queryKey: QueryKey = ["competitor", { id: competitorId }];
      client.setQueryData<Competitor>(queryKey, (competitor) => {
        const currentResult = competitor?.result;

        if (
          currentResult &&
          (!updatedCompetitor.result ||
            currentResult.timestamp > updatedCompetitor.result.timestamp)
        ) {
          return { ...updatedCompetitor, result: currentResult };
        }

        return updatedCompetitor;
      });
    },
  }));
};

export const scrubCompetitorMutation = (competitorId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: async () =>
      ApiClient.getInstance().scrubCompetitor(competitorId),
    onSuccess: (updatedCompetitor) => {
      const queryKey: QueryKey = ["competitor", { id: competitorId }];
      client.setQueryData<Competitor>(queryKey, updatedCompetitor);
    },
  }));
};

export const createCompetitorsMutation = (competitionId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: async (args: CreateCompetitorsArguments) =>
      ApiClient.getInstance().createCompetitors(competitionId, args),
    onSuccess: (newCompetitors) => {
      const queryKey: QueryKey = ["competitors", { competitionId }];
      client.setQueryData<Competitor[]>(queryKey, (oldCompetitors) => {
        return [...(oldCompetitors ?? []), ...newCompetitors];
      });
    },
  }));
};

export const updateCompetitorPublicInfoInQueryCache = (
  queryClient: QueryClient,
  competitorId: number,
  updatedPublicInfo: Pick<
    Competitor,
    | "categoryId"
    | "name"
    | "withdrawnFromFinals"
    | "disqualified"
    | "scrubbedAt"
  >,
) => {
  const queryKey: QueryKey = ["competitor", { id: competitorId }];

  queryClient.setQueryData<Competitor>(queryKey, (competitor) => {
    if (competitor) {
      return { ...competitor, ...updatedPublicInfo };
    }
  });
};

export const updateCompetitorResultInQueryCache = (
  queryClient: QueryClient,
  competitorId: number,
  updatedResult: Result,
) => {
  const queryKey: QueryKey = ["competitor", { id: competitorId }];

  queryClient.setQueryData<Competitor>(queryKey, (competitor) => {
    if (competitor) {
      return {
        ...competitor,
        result: {
          competitorId: updatedResult.competitorId,
          timestamp: updatedResult.timestamp,
          score: updatedResult.score,
          placement: updatedResult.placement,
          rankOrder: updatedResult.rankOrder,
          finalist: updatedResult.finalist,
        },
      };
    }
  });
};

export const updatePointValueInQueryCache = (
  queryClient: QueryClient,
  competitorId: number,
  updatedPointValue: PointValue,
) => {
  const queryKey: QueryKey = ["point-values", { competitorId }];

  queryClient.setQueryData<PointValue[]>(queryKey, (pointValues) => {
    if (pointValues === undefined) {
      return undefined;
    }

    const index = pointValues.findIndex(
      ({ problemId }) => problemId === updatedPointValue.problemId,
    );

    if (index === -1) {
      return [...pointValues, updatedPointValue];
    }

    return pointValues.map((pointValue, currentIndex) => {
      if (currentIndex === index) {
        return updatedPointValue;
      }

      return pointValue;
    });
  });
};
