<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { deleteCompClassMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import type { Snippet } from "svelte";

  type Props = {
    compClassId: number;
    children: Snippet<[{ deleteCompClass: () => void }]>;
  };

  let dialog: WaDialog | undefined = $state();

  let { compClassId, children }: Props = $props();

  const deleteCompClass = $derived(deleteCompClassMutation(compClassId));

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
    $deleteCompClass.mutate(undefined, {
      onError: () => toastError("Failed to delete comp class."),
    });
  };
</script>

{@render children({ deleteCompClass: handleDelete })}

<wa-dialog bind:this={dialog} label="Are you sure?">
  A comp class is deleted permanently and cannot be restored.
  <wa-button slot="footer" appearance="plain" onclick={handleCancel}>
    Cancel</wa-button
  >
  <wa-button slot="footer" variant="danger" onclick={confirmDelete}
    >Remove
    <wa-icon slot="start" name="trash"></wa-icon>
  </wa-button>
</wa-dialog>
