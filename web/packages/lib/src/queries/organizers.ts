import {
  createMutation,
  createQuery,
  useQueryClient,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { Organizer, OrganizerTemplate, User } from "../models";

export const createOrganizerMutation = () => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: async (template: OrganizerTemplate) =>
      ApiClient.getInstance().createOrganizer(template),
    onSuccess: (newOrganizer: Organizer) => {
      client.setQueryData<User>(["self"], (current) => {
        if (!current) {
          return current;
        }

        return {
          ...current,
          organizers: [...(current.organizers ?? []), newOrganizer],
        };
      });
    },
  }));
};

export const getOrganizerQuery = (organizerId: number) =>
  createQuery(() => ({
    queryKey: ["organizer", { id: organizerId }],
    queryFn: async () => ApiClient.getInstance().getOrganizer(organizerId),
  }));
