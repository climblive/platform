import { createQuery } from "@tanstack/svelte-query";
import { ApiClient } from "../Api";

export const getOrganizerQuery = (organizerId: number) =>
  createQuery(() => ({
    queryKey: ["organizer", { id: organizerId }],
    queryFn: async () => ApiClient.getInstance().getOrganizer(organizerId),
  }));
