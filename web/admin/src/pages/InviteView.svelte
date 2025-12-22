<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/copy-button/copy-button.js";
  import { ApiClient } from "@climblive/lib";
  import type { OrganizerInviteID } from "@climblive/lib/models";
  import {
    deleteOrganizerInviteMutation,
    getOrganizerInviteQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  interface Props {
    inviteId: OrganizerInviteID;
  }

  const { inviteId }: Props = $props();

  const inviteQuery = $derived(getOrganizerInviteQuery(inviteId));
  const deleteInvite = $derived(deleteOrganizerInviteMutation(inviteId));

  const invite = $derived(inviteQuery.data);

  const handleAccept = async () => {
    if (!invite) {
      return;
    }

    try {
      await ApiClient.getInstance().acceptOrganizerInvite(invite.id);

      navigate(`./organizers/${invite.organizerId}`);
    } catch {
      toastError("Failed to accept invite.");
    }
  };

  const handleDecline = () => {
    deleteInvite.mutate(undefined, {
      onSuccess: () => {
        navigate(`./`);
      },
      onError: () => {
        toastError("Failed to decline invite.");
      },
    });
  };
</script>

{#if invite}
  <p>
    You have been invited to be a member of <strong
      >{invite.organizerName}</strong
    >.
  </p>
  <section>
    <wa-button
      variant="danger"
      appearance="outlined"
      onclick={handleDecline}
      loading={deleteInvite.isPending}
      >Decline
      <wa-icon slot="start" name="trash"></wa-icon>
    </wa-button>
    <wa-button
      variant="success"
      appearance="outlined filled"
      onclick={handleAccept}
      loading={false}
      >Accept
      <wa-icon slot="start" name="check"></wa-icon>
    </wa-button>
  </section>
{/if}

<style>
  section {
    display: flex;
    align-items: start;
    gap: var(--wa-space-m);
  }
</style>
