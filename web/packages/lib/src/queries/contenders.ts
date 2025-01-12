import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Contender, ContenderPatch } from "../models";
import { HOUR } from "./constants";

export const getContenderQuery = (contenderId: number) =>
  createQuery({
    queryKey: ["contender", { id: contenderId }],
    queryFn: async () => ApiClient.getInstance().getContender(contenderId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 0,
    refetchOnWindowFocus: true,
  });

export const updateContenderMutation = (contenderId: number) => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: async (patch: ContenderPatch) =>
      ApiClient.getInstance().updateContender(contenderId, patch),
    onSuccess: (updatedContender) => {
      const queryKey: QueryKey = ["contender", { id: contenderId }];
      client.setQueryData<Contender>(queryKey, updatedContender);
    },
  });
};
