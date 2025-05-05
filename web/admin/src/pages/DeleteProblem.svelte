<script lang="ts">
  import { deleteProblemMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import type { SlDialog } from "@shoelace-style/shoelace";
  import type { Snippet } from "svelte";

  type Props = {
    problemID: number;
    children: Snippet<[{ onDelete: () => void }]>;
  };

  let dialog: SlDialog | undefined = $state();

  let { problemID, children }: Props = $props();

  const deleteProblem = deleteProblemMutation(problemID);

  const handleDelete = async () => {
    console.log("Deleting problem with ID:", problemID);

    dialog?.show();
  };

  const confirmDelete = () => {
    $deleteProblem.mutate(undefined, {
      onError: () => toastError("Failed to delete problem."),
    });
  };
</script>

{@render children({ onDelete: handleDelete })}

<sl-dialog bind:this={dialog} no-header>
  <p>
    <strong>Are you sure?</strong>
  </p>
  <p>A problem is deleted permanently and cannot be restored.</p>
  <sl-button slot="footer" variant="text" onclick={() => dialog?.hide()}
    >Cancel</sl-button
  >
  <sl-button slot="footer" variant="danger" onclick={confirmDelete}
    >Remove
    <sl-icon slot="prefix" name="trash"></sl-icon>
  </sl-button>
</sl-dialog>
