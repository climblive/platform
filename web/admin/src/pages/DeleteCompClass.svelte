<script lang="ts">
  import { deleteCompClassMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import type { SlDialog } from "@shoelace-style/shoelace";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import "@shoelace-style/shoelace/dist/components/dialog/dialog.js";
  import "@shoelace-style/shoelace/dist/components/icon/icon.js";
  import type { Snippet } from "svelte";

  type Props = {
    compClassId: number;
    children: Snippet<[{ deleteCompClass: () => void }]>;
  };

  let dialog: SlDialog | undefined = $state();

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

<sl-dialog bind:this={dialog} no-header>
  <p>
    <strong>Are you sure?</strong>
  </p>
  <p>A comp class is deleted permanently and cannot be restored.</p>
  <sl-button slot="footer" variant="text" onclick={() => dialog?.hide()}
    >Cancel</sl-button
  >
  <sl-button slot="footer" variant="danger" onclick={confirmDelete}
    >Remove
    <sl-icon slot="prefix" name="trash"></sl-icon>
  </sl-button>
</sl-dialog>
