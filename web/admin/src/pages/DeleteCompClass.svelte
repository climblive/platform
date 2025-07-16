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
    dialog?.show();
  };

  const confirmDelete = () => {
    $deleteCompClass.mutate(undefined, {
      onError: () => toastError("Failed to delete comp class."),
    });
  };
</script>

{@render children({ deleteCompClass: handleDelete })}

<wa-dialog bind:this={dialog} no-header>
  <p>
    <strong>Are you sure?</strong>
  </p>
  <p>A comp class is deleted permanently and cannot be restored.</p>
  <wa-button slot="footer" variant="text" onclick={() => dialog?.hide()}
    >Cancel</wa-button
  >
  <wa-button slot="footer" variant="danger" onclick={confirmDelete}
    >Remove
    <wa-icon slot="prefix" name="trash"></wa-icon>
  </wa-button>
</wa-dialog>
