<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { createUnlockRequestMutation } from "@climblive/lib/queries";
  import { toastError, toastSuccess } from "@climblive/lib/utils";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let dialog: WaDialog | undefined = $state();

  const createRequest = $derived(createUnlockRequestMutation());

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

  const handleRequest = () => {
    createRequest.mutate(
      { contestId },
      {
        onSuccess: (request) => {
          if (request.status === "approved") {
            toastSuccess(
              "Full capacity unlocked! You can now create up to 500 tickets.",
            );
          } else {
            toastSuccess(
              "Request submitted! You'll be notified once it's reviewed.",
            );
          }
          close();
        },
        onError: () => toastError("Failed to submit request."),
      },
    );
  };
</script>

<wa-dialog bind:this={dialog} label="Request full capacity">
  <div class="dialog-content">
    <wa-callout variant="warning">
      <wa-icon slot="icon" name="triangle-exclamation"></wa-icon>
      Help us keep things organized by only requesting this for real contests.
    </wa-callout>

    <p>
      Evaluation mode limits contests to 10 contenders for testing purposes. By
      requesting full capacity you'll be able to host up to 500 contenders once
      approved.
    </p>

    <p>
      If you've been approved before, your request will be automatically
      approved.
    </p>
  </div>

  <wa-button slot="footer" appearance="plain" onclick={close}>Cancel</wa-button>
  <wa-button
    slot="footer"
    size="small"
    variant="success"
    appearance="accent"
    loading={createRequest.isPending}
    onclick={handleRequest}
  >
    <wa-icon slot="start" name="paper-plane"></wa-icon>
    Submit request
  </wa-button>
</wa-dialog>

<style>
  .dialog-content {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }
</style>
