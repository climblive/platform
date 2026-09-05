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

export const putTickMutation = (contenderId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (tick: Omit<Tick, "timestamp">) =>
      ApiClient.getInstance().putTick(contenderId, tick),
    onSuccess: (updatedTick) => {
      updateTickInQueryCache(client, contenderId, updatedTick);
    },
  }));
};

export const deleteTickMutation = (contenderId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (problemId: number) =>
      ApiClient.getInstance().deleteTick(contenderId, problemId),
    onSuccess: (...args) => {
      const [, problemId] = args;

      removeTickFromQueryCache(client, contenderId, problemId);
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
    const predicate = ({ problemId }: Tick) =>
      problemId === updatedTick.problemId;

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
  contenderId: number,
  problemId: number,
) => {
  const queryKey: QueryKey = ["ticks", { contenderId }];

  queryClient.setQueryData<Tick[]>(queryKey, (oldTicks) => {
    const predicate = (tick: Tick) => tick.problemId !== problemId;

    return oldTicks ? oldTicks.filter(predicate) : undefined;
  });
};
