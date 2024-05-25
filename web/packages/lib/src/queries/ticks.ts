import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Tick } from "../models/tick";
import { HOUR } from "./constants";

export const getTicksQuery = (contenderId: number) =>
  createQuery({
    queryKey: ["ticks", { contenderId }],
    queryFn: async () => ApiClient.getInstance().getTicks(contenderId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
  });

export const createTickMutation = (contenderId: number) => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: (tick: Omit<Tick, "id">) =>
      ApiClient.getInstance().createTick(contenderId, tick),
    onSuccess: (newTick) => {
      const queryKey: QueryKey = ["ticks", { contenderId }];

      client.setQueryData<Tick[]>(queryKey, (oldTicks) =>
        oldTicks ? [...oldTicks, newTick] : [newTick]
      );
    },
  });
};

export const deleteTickMutation = () => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: (tickId: number) => ApiClient.getInstance().deleteTick(tickId),
    onSuccess: (...args) => {
      const [, variables] = args;
      const queryKey = ["ticks"];
      client.setQueriesData<Tick[]>(
        {
          queryKey,
          exact: false,
        },
        (oldTicks) =>
          oldTicks ? oldTicks.filter(({ id }) => id !== variables) : undefined
      );
    },
  });
};
