import { createQuery } from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { User } from "../models";

export const getSelfQuery = () =>
  createQuery(() => ({
    queryKey: ["self"],
    queryFn: async () => ApiClient.getInstance().getSelf(),
  });

export const getUsersByOrganizerQuery = (
  organizerId: number,
  options?: Partial<Parameters<typeof createQuery<User[]>>[0]>,
) =>
  createQuery({
    ...options,
    queryKey: ["users", { organizerId }],
    queryFn: async () =>
      ApiClient.getInstance().getUsersByOrganizer(organizerId),
  });
