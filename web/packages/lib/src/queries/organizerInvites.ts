import {
  createMutation,
  createQuery,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import type { OrganizerInvite, OrganizerInviteID } from "../models";

export const getOrganizerInvitesQuery = (organizerId: number) =>
  createQuery(() => ({
    queryKey: ["organizer-invites", { organizerId }],
    queryFn: async () =>
      ApiClient.getInstance().getOrganizerInvites(organizerId),
  }));

export const getOrganizerInviteQuery = (inviteId: OrganizerInviteID) =>
  createQuery(() => ({
    queryKey: ["organizer-invite", { id: inviteId }],
    queryFn: async () => ApiClient.getInstance().getOrganizerInvite(inviteId),
  }));

export const createOrganizerInviteMutation = (organizerId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: () =>
      ApiClient.getInstance().createOrganizerInvite(organizerId),
    onSuccess: (newInvite) => {
      let queryKey: QueryKey = ["organizer-invites", { organizerId }];

      client.setQueryData<OrganizerInvite[]>(queryKey, (oldInvites) => {
        return [...(oldInvites ?? []), newInvite];
      });

      queryKey = ["organizer-invite", { id: newInvite.id }];

      client.setQueryData<OrganizerInvite>(queryKey, newInvite);
    },
  }));
};

export const deleteOrganizerInviteMutation = (inviteId: OrganizerInviteID) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: () => ApiClient.getInstance().deleteOrganizerInvite(inviteId),
    onSuccess: () => {
      let queryKey: QueryKey = ["organizer-invites"];

      client.setQueriesData<OrganizerInvite[]>(
        { queryKey, exact: false },
        (oldInvites) => {
          if (oldInvites === undefined) {
            return undefined;
          }

          return oldInvites.filter(({ id }) => id !== inviteId);
        },
      );

      queryKey = ["organizer-invite", { id: inviteId }];

      client.removeQueries({ queryKey });
    },
  }));
};
