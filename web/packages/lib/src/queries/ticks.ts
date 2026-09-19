import {
  createMutation,
  createQuery,
  QueryClient,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Tick } from "../models";

export const getTicksByContenderQuery = (
  contenderId: number,
  options?: Partial<Parameters<typeof createQuery<Tick[]>>[0]>,
) =>
  createQuery<Tick[]>(() => ({
    ...options,
    queryKey: ["ticks", { contenderId }],
    queryFn: async () =>
      ApiClient.getInstance().getTicksByContender(contenderId),
    refetchOnWindowFocus: "always",
  }));

export const getTicksByContestQuery = (contestId: number) =>
  createQuery(() => ({
    queryKey: ["ticks", { contestId }],
    queryFn: async () => ApiClient.getInstance().getTicksByContest(contestId),
  }));

export const putTickMutation = (contenderId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (tick: Omit<Tick, "id" | "timestamp">) =>
      ApiClient.getInstance().putTick(contenderId, tick),
    onSuccess: (updatedTick) => {
      updateTickInQueryCache(client, contenderId, updatedTick);
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
