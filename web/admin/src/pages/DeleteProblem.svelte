<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { deleteProblemMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import type { Snippet } from "svelte";

  type Props = {
    problemId: number;
    children: Snippet<[{ deleteProblem: () => void }]>;
  };

  let dialog: WaDialog | undefined = $state();

  let { problemId, children }: Props = $props();

  const deleteProblem = $derived(deleteProblemMutation(problemId));

  const handleDelete = async () => {
    dialog?.setAttribute("open", "true");
  };

  const confirmDelete = () => {
    $deleteProblem.mutate(undefined, {
      onError: () => toastError("Failed to delete problem."),
    });
  };
</script>

{@render children({ deleteProblem: handleDelete })}

<wa-dialog bind:this={dialog} no-header>
  <p>
    <strong>Are you sure?</strong>
  </p>
  <p>A problem is deleted permanently and cannot be restored.</p>
  <wa-button
    slot="footer"
    variant="text"
    onclick={() => dialog?.setAttribute("open", "false")}>Cancel</wa-button
  >
  <wa-button slot="footer" variant="danger" onclick={confirmDelete}
    >Remove
    <wa-icon slot="prefix" name="trash"></wa-icon>
  </wa-button>
</wa-dialog>
