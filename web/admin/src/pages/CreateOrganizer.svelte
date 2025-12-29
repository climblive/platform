<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import type WaInput from "@awesome.me/webawesome/dist/components/input/input.js";
  import type { Organizer } from "@climblive/lib/models";
  import { createOrganizerMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";
  import type { Snippet } from "svelte";

  type Props = {
    children: Snippet<[{ createOrganizer: () => void }]>;
  };

  let dialog: WaDialog | undefined = $state();
  let nameInput: WaInput | undefined = $state();

  const { children }: Props = $props();

  const createOrganizer = $derived(createOrganizerMutation());

  const handleOpen = () => {
    if (dialog) {
      dialog.open = true;
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

    if (!name) {
      return;
    }

    createOrganizer.mutate(
      { name },
      {
        onSuccess: (organizer: Organizer) => {
          if (dialog) {
            dialog.open = false;
          }
          navigate(`/admin/organizers/${organizer.id}/contests`);
        },
        onError: () => toastError("Failed to create organizer."),
      },
    );
  };
</script>

{@render children({ createOrganizer: handleOpen })}

<wa-dialog bind:this={dialog} label="Create new organizer">
  <form onsubmit={handleSubmit}>
    <wa-input
      bind:this={nameInput}
      label="Organizer name"
      required
      autofocus
    ></wa-input>

    <div class="controls">
      <wa-button appearance="plain" onclick={handleCancel} type="button">
        Cancel
      </wa-button>
      <wa-button
        variant="neutral"
        appearance="accent"
        type="submit"
        loading={createOrganizer.isPending}
      >
        Create
      </wa-button>
    </div>
  </form>
</wa-dialog>

<style>
  wa-dialog {
    white-space: normal;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }

  .controls {
    display: flex;
    gap: var(--wa-space-s);
    justify-content: flex-end;
  }
</style>
