import {
  createMutation,
  createQuery,
  QueryClient,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Tick } from "../models";
import { HOUR } from "./constants";

export const getTicksByContenderQuery = (
  contenderId: number,
  options?: Partial<Parameters<typeof createQuery<Tick[]>>[0]>,
) =>
  createQuery<Tick[]>(() => ({
    ...options,
    queryKey: ["ticks", { contenderId }],
    queryFn: async () =>
      ApiClient.getInstance().getTicksByContender(contenderId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 0,
    refetchOnWindowFocus: true,
  }));

export const getTicksByContestQuery = (contestId: number) =>
  createQuery(() => ({
    queryKey: ["ticks", { contestId }],
    queryFn: async () => ApiClient.getInstance().getTicksByContest(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 0,
    refetchOnWindowFocus: true,
  }));

export const createTickMutation = (contenderId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (tick: Omit<Tick, "id" | "timestamp">) =>
      ApiClient.getInstance().createTick(contenderId, tick),
    onSuccess: (newTick) => {
      updateTickInQueryCache(client, contenderId, newTick);
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

export const updateTickInQueryCache = (
  queryClient: QueryClient,
  contenderId: number,
  updatedTick: Tick,
) => {
  const queryKey: QueryKey = ["ticks", { contenderId }];

  queryClient.setQueryData<Tick[]>(queryKey, (oldTicks) => {
    const predicate = ({ id }: Tick) => id === updatedTick.id;

    const found = (oldTicks ?? []).findIndex(predicate) !== -1;

    if (found) {
      return (oldTicks ?? []).map((oldTick) =>
        predicate(oldTick) ? updatedTick : oldTick,
      );
    } else {
      return [...(oldTicks ?? []), updatedTick];
    }
  });
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
