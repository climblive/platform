import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Raffle, RaffleWinner } from "../models";
import { HOUR } from "./constants";

export const getRaffleQuery = (raffleId: number) =>
  createQuery({
    queryKey: ["raffle", { id: raffleId }],
    queryFn: async () => ApiClient.getInstance().getRaffle(raffleId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
  });

export const getRafflesQuery = (contestId: number) =>
  createQuery({
    queryKey: ["raffles", { contestId }],
    queryFn: async () => ApiClient.getInstance().getRaffles(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
  });

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

export const drawRaffleWinnerMutation = (raffleId: number) => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: () => ApiClient.getInstance().drawRaffleWinner(raffleId),
    onSuccess: (newWinner) => {
      const queryKey: QueryKey = ["raffle-winners", { raffleId }];

      client.setQueryData<RaffleWinner[]>(queryKey, (oldWinners) => {
        return [...(oldWinners ?? []), newWinner];
      });
    },
  });
};
