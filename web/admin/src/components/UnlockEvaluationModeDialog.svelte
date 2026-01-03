<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { patchContestMutation } from "@climblive/lib/queries";
  import { toastError, toastSuccess } from "@climblive/lib/utils";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let dialog: WaDialog | undefined = $state();

  const patchContest = $derived(patchContestMutation(contestId));

  export const open = () => {
    if (dialog) {
      dialog.open = true;
    }
  };

  const close = () => {
    if (dialog) {
      dialog.open = false;
    }
  };

  const handleUnlock = () => {
    patchContest.mutate(
      { evaluationMode: false },
      {
        onSuccess: () => {
          toastSuccess(
            "Evaluation mode unlocked. You can now create up to 500 tickets.",
          );
          close();
        },
        onError: () => toastError("Failed to unlock evaluation mode."),
      },
    );
  };
</script>

<wa-dialog bind:this={dialog} label="Unlock full capacity">
  <div class="dialog-content">
    <wa-callout variant="warning">
      <wa-icon slot="icon" name="triangle-exclamation"></wa-icon>
      Help us keep things organized by only unlocking this for real contests.
    </wa-callout>

    <p>
      Evaluation mode limits contests to 10 contenders for testing purposes. By
      unlocking the full capacity you'll be able to host up to 500 contenders.
    </p>
  </div>

  <wa-button slot="footer" appearance="plain" onclick={close}>Cancel</wa-button>
  <wa-button
    slot="footer"
    size="small"
    variant="success"
    appearance="accent"
    loading={patchContest.isPending}
    onclick={handleUnlock}
  >
    <wa-icon slot="start" name="lock-open"></wa-icon>
    Unlock
  </wa-button>
</wa-dialog>

<style>
  .dialog-content {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }
</style>
