<script lang="ts">
  import { type WaSelectEvent } from "@awesome.me/webawesome";
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/dropdown-item/dropdown-item.js";
  import WaDropdownItem from "@awesome.me/webawesome/dist/components/dropdown-item/dropdown-item.js";
  import "@awesome.me/webawesome/dist/components/dropdown/dropdown.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { stopScoreEngineMutation } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";

  type Props = {
    instanceId: string;
  };

  const { instanceId }: Props = $props();

  let dialog: WaDialog | undefined = $state();

  const stopScoreEngine = stopScoreEngineMutation();

  const handleCancel = () => {
    if (dialog) {
      dialog.open = false;
    }
  };

  const handleStop = () => {
    if (dialog) {
      dialog.open = true;
    }
  };

  const confirmStop = () => {
    stopScoreEngine.mutate(instanceId, {
      onSettled: () => {
        if (dialog) {
          dialog.open = false;
        }
      },
      onError: () => toastUnexpectedError("Failed to stop score engine"),
    });
  };
</script>

<wa-dropdown
  onwa-select={(event: WaSelectEvent) => {
    if ((event.detail.item as WaDropdownItem).value === "delete") {
      handleStop();
    }
  }}
>
  <wa-button slot="trigger" size="s" appearance="plain">
    <wa-icon name="ellipsis-vertical" label="Actions"></wa-icon>
  </wa-button>
  <wa-dropdown-item value="delete" variant="danger">
    <wa-icon slot="icon" name="stop"></wa-icon>
    Stop
  </wa-dropdown-item>
</wa-dropdown>

<wa-dialog bind:this={dialog} label="Stop score engine">
  Stopping the score engine may interrupt live scoring.
  <wa-button slot="footer" appearance="plain" onclick={handleCancel}
    >Cancel</wa-button
  >
  <wa-button
    slot="footer"
    variant="danger"
    onclick={confirmStop}
    loading={stopScoreEngine.isPending}
  >
    Proceed
    <wa-icon slot="stop" name="trash"></wa-icon>
  </wa-button>
</wa-dialog>

<style>
  wa-dialog {
    white-space: normal;
  }
</style>
