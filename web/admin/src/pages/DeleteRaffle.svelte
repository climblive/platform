<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { deleteRaffleMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import type { Snippet } from "svelte";

  type Props = {
    raffleId: number;
    children: Snippet<[{ deleteRaffle: () => void }]>;
  };

  let dialog: WaDialog | undefined = $state();

  let { raffleId, children }: Props = $props();

  const deleteRaffle = $derived(deleteRaffleMutation(raffleId));

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
    deleteRaffle.mutate(undefined, {
      onError: () => toastError("Failed to delete raffle."),
    });
  };
</script>

{@render children({ deleteRaffle: handleDelete })}

<wa-dialog bind:this={dialog} label="Delete raffle">
  Deleting a raffle will permanently remove it and it cannot be restored.
  <wa-button slot="footer" appearance="plain" onclick={handleCancel}>
    Cancel</wa-button
  >
  <wa-button
    slot="footer"
    variant="danger"
    onclick={confirmDelete}
    loading={deleteRaffle.isPending}
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
