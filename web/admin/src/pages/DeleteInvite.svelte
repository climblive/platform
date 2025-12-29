<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import type { OrganizerInviteID } from "@climblive/lib/models";
  import { deleteOrganizerInviteMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import type { Snippet } from "svelte";

  type Props = {
    inviteId: OrganizerInviteID;
    children: Snippet<[{ deleteInvite: () => void }]>;
  };

  let dialog: WaDialog | undefined = $state();

  const { inviteId, children }: Props = $props();

  const deleteInvite = $derived(deleteOrganizerInviteMutation(inviteId));

  const handleDelete = async () => {
    if (dialog) {
      dialog.open = true;
    }
  };

  const handleCancel = () => {
    if (dialog) {
      dialog.open = false;
    }
  };

  const confirmDelete = () => {
    deleteInvite.mutate(undefined, {
      onError: () => toastError("Failed to delete invite."),
    });
  };
</script>

{@render children({ deleteInvite: handleDelete })}

<wa-dialog bind:this={dialog} label="Delete invite">
  Deleting an invite will permanently remove it and it cannot be restored.
  <wa-button slot="footer" appearance="plain" onclick={handleCancel}>
    Cancel</wa-button
  >
  <wa-button
    slot="footer"
    variant="danger"
    onclick={confirmDelete}
    loading={deleteInvite.isPending}
  >
    Remove
    <wa-icon slot="start" name="trash"></wa-icon>
  </wa-button>
</wa-dialog>

<style>
  wa-dialog {
    white-space: normal;
  }
</style>
