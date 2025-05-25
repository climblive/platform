import {
  createMutation,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Raffle } from "../models";

export const createRaffleMutation = (contestId: number) => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: () => ApiClient.getInstance().createRaffle(contestId),
    onSuccess: (newRaffle) => {
      let queryKey: QueryKey = ["raffles", { contestId }];

      client.setQueryData<Raffle[]>(queryKey, (oldRaffles) => {
        return [...(oldRaffles ?? []), newRaffle];
      });

      queryKey = ["raffle", { id: newRaffle.id }];

      client.setQueryData<Raffle>(queryKey, newRaffle);
    },
  });
};
