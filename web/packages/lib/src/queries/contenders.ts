import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Contender, ContenderPatch } from "../models";
import type { CreateContendersArguments } from "../models/rest";
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

export const getContendersByContestQuery = (contestId: number) =>
  createQuery({
    queryKey: ["contenders", { contestId }],
    queryFn: async () =>
      ApiClient.getInstance().getContendersByContest(contestId),
    retry: false,
    gcTime: 12 * HOUR,
    staleTime: 0,
  });

export const patchContenderMutation = (contenderId: number) => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: async (patch: ContenderPatch) =>
      ApiClient.getInstance().patchContender(contenderId, patch),
    onSuccess: (updatedContender) => {
      const queryKey: QueryKey = ["contender", { id: contenderId }];
      client.setQueryData<Contender>(queryKey, updatedContender);
    },
  });
};

export const createContendersMutation = (contestId: number) => {
  const client = useQueryClient();

  return createMutation({
    mutationFn: async (args: CreateContendersArguments) =>
      ApiClient.getInstance().createContenders(contestId, args),
    onSuccess: (newContenders) => {
      const queryKey: QueryKey = ["contenders", { contestId }];
      client.setQueryData<Contender[]>(queryKey, (oldContenders) => {
        return [...(oldContenders ?? []), ...newContenders];
      });
    },
  });
};
