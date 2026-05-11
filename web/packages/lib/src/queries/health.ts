import { createQuery } from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import { SECOND } from "./constants";

export const getHealthQuery = () =>
  createQuery(() => ({
    queryKey: ["health"],
    queryFn: async () => ApiClient.getInstance().getHealth(),
    refetchInterval: 30 * SECOND,
    retry: false,
  }));
