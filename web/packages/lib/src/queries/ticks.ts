import {
  createMutation,
  createQuery,
  QueryClient,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Tick } from "../models";

export const getTicksByCompetitorQuery = (
  competitorId: number,
  options?: Partial<Parameters<typeof createQuery<Tick[]>>[0]>,
) =>
  createQuery<Tick[]>(() => ({
    ...options,
    queryKey: ["ticks", { competitorId }],
    queryFn: async () =>
      ApiClient.getInstance().getTicksByCompetitor(competitorId),
    refetchOnWindowFocus: "always",
  }));

export const getTicksByCompetitionQuery = (competitionId: number) =>
  createQuery(() => ({
    queryKey: ["ticks", { competitionId }],
    queryFn: async () =>
      ApiClient.getInstance().getTicksByCompetition(competitionId),
  }));

export const putTickMutation = (competitorId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (tick: Omit<Tick, "id" | "timestamp">) =>
      ApiClient.getInstance().putTick(competitorId, tick),
    onSuccess: (updatedTick) => {
      updateTickInQueryCache(client, competitorId, updatedTick);
    },
  }));
};

export const deleteTickMutation = () => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (tickId: number) => ApiClient.getInstance().deleteTick(tickId),
    onSuccess: (...args) => {
      const [, tickId] = args;

      removeTickFromQueryCache(client, tickId);
    },
  }));
};

const updateTick = (oldTicks: Tick[] | undefined, updatedTick: Tick) => {
  const predicate = ({ id }: Tick) => id === updatedTick.id;

  const existingTick = (oldTicks ?? []).find(predicate);

  if (existingTick && existingTick.revision > updatedTick.revision) {
    return oldTicks;
  }

  if (existingTick) {
    return (oldTicks ?? []).map((oldTick) =>
      predicate(oldTick) ? updatedTick : oldTick,
    );
  } else {
    return [...(oldTicks ?? []), updatedTick];
  }
};

export const updateTickInQueryCache = (
  queryClient: QueryClient,
  competitorId: number,
  updatedTick: Tick,
) => {
  const queryKey: QueryKey = ["ticks", { competitorId }];

  queryClient.setQueryData<Tick[]>(queryKey, (oldTicks) =>
    updateTick(oldTicks, updatedTick),
  );
};

export const updateTickInCompetitionQueryCache = (
  queryClient: QueryClient,
  competitionId: number,
  updatedTick: Tick,
) => {
  const queryKey: QueryKey = ["ticks", { competitionId }];

  queryClient.setQueryData<Tick[]>(queryKey, (oldTicks) =>
    updateTick(oldTicks, updatedTick),
  );
};

export const removeTickFromQueryCache = (
  queryClient: QueryClient,
  tickId: number,
) => {
  const queryKey = ["ticks"];

  queryClient.setQueriesData<Tick[]>(
    {
      queryKey,
      exact: false,
    },
    (oldTicks) => {
      const predicate = ({ id }: Tick) => id !== tickId;

      return oldTicks ? oldTicks.filter(predicate) : undefined;
    },
  );
};
