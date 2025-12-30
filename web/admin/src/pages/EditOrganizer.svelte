<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import type WaInput from "@awesome.me/webawesome/dist/components/input/input.js";
  import type { Organizer } from "@climblive/lib/models";
  import { patchOrganizerMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import type { Snippet } from "svelte";

  type Props = {
    organizer: Organizer;
    children: Snippet<[{ editOrganizer: () => void }]>;
  };

  let dialog: WaDialog | undefined = $state();
  let nameInput: WaInput | undefined = $state();

  const { organizer, children }: Props = $props();

  const patchOrganizer = $derived(patchOrganizerMutation(organizer.id));

  const handleOpen = () => {
    if (dialog) {
      dialog.open = true;

      if (nameInput) {
        nameInput.value = organizer.name;
      }
    }
  };

  const handleCancel = () => {
    if (dialog) {
      dialog.open = false;
    }
  };

  const handleSubmit = (e: SubmitEvent) => {
    e.preventDefault();

    const name = nameInput?.value?.trim();

    if (!name || name === organizer.name) {
      handleCancel();
      return;
    }

    patchOrganizer.mutate(
      { name },
      {
        onSuccess: () => {
          handleCancel();
        },
        onError: () => toastError("Failed to update organizer name."),
      },
    );
  };
</script>

<div>
  {@render children({ editOrganizer: handleOpen })}

  <wa-dialog bind:this={dialog} label="Edit organizer">
    <form onsubmit={handleSubmit}>
      <wa-input bind:this={nameInput} label="Organizer name" required
      ></wa-input>

      <div class="controls">
        <wa-button appearance="plain" onclick={handleCancel} type="button">
          Cancel
        </wa-button>
        <wa-button
          variant="neutral"
          appearance="accent"
          type="submit"
          loading={patchOrganizer.isPending}
        >
          Save
        </wa-button>
      </div>
    </form>
  </wa-dialog>
</div>

<style>
  .controls {
    display: flex;
    justify-content: flex-end;
    gap: var(--wa-space-xs);
    margin-top: var(--wa-space-m);
  }
</style>
