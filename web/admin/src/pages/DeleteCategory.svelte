<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { deleteCategoryMutation } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import type { Snippet } from "svelte";

  type Props = {
    categoryId: number;
    children: Snippet<[{ deleteCategory: () => void }]>;
  };

  let dialog: WaDialog | undefined = $state();

  let { categoryId, children }: Props = $props();

  const deleteCategory = $derived(deleteCategoryMutation(categoryId));

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
    deleteCategory.mutate(undefined, {
      onError: () => toastUnexpectedError("Failed to delete category."),
    });
  };
</script>

{@render children({ deleteCategory: handleDelete })}

<wa-dialog bind:this={dialog} label="Delete category">
  Are you sure you want to delete this category? This action is permanent and
  cannot be undone.
  <wa-button slot="footer" appearance="plain" onclick={handleCancel}>
    Cancel</wa-button
  >
  <wa-button
    slot="footer"
    variant="danger"
    onclick={confirmDelete}
    loading={deleteCategory.isPending}
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
