import { createMutation, useQueryClient, type QueryKey } from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Raffle, RaffleTemplate } from "../models";

export const createRaffleMutation = (contestId: number) => {
    const client = useQueryClient();

    return createMutation({
        mutationFn: (template: RaffleTemplate) =>
            ApiClient.getInstance().createRaffle(contestId, template),
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