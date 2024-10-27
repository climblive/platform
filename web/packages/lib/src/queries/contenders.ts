import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Contender } from "../models/contender";
import { HOUR } from "./constants";

export const getContenderQuery = (contenderId: number) =>
  createQuery({
    queryKey: ["contender", { id: contenderId }],
    queryFn: async () => ApiClient.getInstance().getContender(contenderId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 12 * HOUR,
    refetchOnWindowFocus: true,
  });

export const updateContenderMutation = (contenderId: number) => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: async (contender: Contender) =>
      ApiClient.getInstance().updateContender(contenderId, contender),
    onSuccess: (updatedContender) => {
      const queryKey: QueryKey = ["contender", { id: contenderId }];
      client.setQueryData<Contender>(queryKey, updatedContender);
    },
  });
};
