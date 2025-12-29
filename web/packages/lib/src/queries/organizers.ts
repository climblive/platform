import {
  createMutation,
  createQuery,
  useQueryClient,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { OrganizerTemplate } from "../models";

export const createOrganizerMutation = () => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: async (template: OrganizerTemplate) =>
      ApiClient.getInstance().createOrganizer(template),
    onSuccess: () => {
      client.invalidateQueries({ queryKey: ["user", "self"] });
    },
  }));
};

export const getOrganizerQuery = (organizerId: number) =>
  createQuery(() => ({
    queryKey: ["organizer", { id: organizerId }],
    queryFn: async () => ApiClient.getInstance().getOrganizer(organizerId),
  }));
