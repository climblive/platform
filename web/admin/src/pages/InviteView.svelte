<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import { ApiClient } from "@climblive/lib";
  import type { OrganizerInviteID } from "@climblive/lib/models";
  import {
    deleteOrganizerInviteMutation,
    getOrganizerInviteQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { isAfter } from "date-fns";
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

      navigate(`./organizers/${invite.organizerId}/contests`);
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
  {#if isAfter(new Date(), invite.expiresAt)}
    <wa-callout variant="danger">
      <wa-icon slot="icon" name="circle-exclamation"></wa-icon>
      This invite has expired and can no longer be accepted.
    </wa-callout>
  {:else}
    <p>
      You have been invited to be a co-organizer of <strong
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
        appearance="filled-outlined"
        onclick={handleAccept}
        >Accept
        <wa-icon slot="start" name="check"></wa-icon>
      </wa-button>
    </section>
  {/if}
{/if}

<style>
  section {
    display: flex;
    align-items: start;
    gap: var(--wa-space-xs);
  }
</style>
