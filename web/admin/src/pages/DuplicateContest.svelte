<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { duplicateContestMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  type Props = {
    contestId: number;
  };

  let dialog: WaDialog | undefined = $state();

  let { contestId }: Props = $props();

  const duplicateContest = $derived(duplicateContestMutation(contestId));

  const handleDuplication = async () => {
    if (dialog) {
      dialog.open = true;
    }
  };

  const handleCancel = () => {
    if (dialog) {
      dialog.open = false;
    }
  };

  const confirmDuplication = () => {
    $duplicateContest.mutate(undefined, {
      onSuccess: (duplicate) => {
        handleCancel();
        navigate(`/admin/contests/${duplicate.id}`);
      },
      onError: () => {
        toastError("Failed to duplicate contest.");
      },
    });
  };
</script>

<div class="actions">
  <wa-button onclick={handleDuplication} appearance="outlined"
    >Duplicate
    <wa-icon name="copy" slot="start"></wa-icon>
  </wa-button>
</div>

<wa-dialog bind:this={dialog} label="Duplicate contest">
  You are about to create a copy of the contest. Everything except tickets,
  results and raffles will be copied.
  <wa-button slot="footer" appearance="plain" onclick={handleCancel}>
    Cancel</wa-button
  >
  <wa-button
    slot="footer"
    variant="warning"
    onclick={confirmDuplication}
    loading={$duplicateContest.isPending}
  >
    Duplicate
    <wa-icon slot="start" name="copy"></wa-icon>
  </wa-button>
</wa-dialog>
