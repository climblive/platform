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

export const getTicksQuery = (contenderId: number) =>
  createQuery({
    queryKey: ["ticks", { contenderId }],
    queryFn: async () => ApiClient.getInstance().getTicks(contenderId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 0,
    refetchOnWindowFocus: true,
  });

export const createTickMutation = (contenderId: number) => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: (tick: Omit<Tick, "id" | "timestamp">) =>
      ApiClient.getInstance().createTick(contenderId, tick),
    onSuccess: (newTick) => {
      updateTickInCache(client, contenderId, newTick);
    },
  });
};

export const deleteTickMutation = () => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: (tickId: number) => ApiClient.getInstance().deleteTick(tickId),
    onSuccess: (...args) => {
      const [, tickId] = args;

      removeTickFromCache(client, tickId);
    },
  });
};

export const updateTickInCache = (
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

export const removeTickFromCache = (
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
