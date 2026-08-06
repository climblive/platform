<script lang="ts">
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { stopScoreEngineMutation } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";

  type Props = {
    instanceId: string;
  };

  const { instanceId }: Props = $props();

  const stopScoreEngine = stopScoreEngineMutation();

  let confirmStopEngineId = $state<string | undefined>(undefined);

  const handleCancel = () => {
    confirmStopEngineId = undefined;
  };

  const handleStop = () => {
    if (confirmStopEngineId !== instanceId) {
      confirmStopEngineId = instanceId;
      return;
    }

    stopScoreEngine.mutate(instanceId, {
      onSettled: () => {
        confirmStopEngineId = undefined;
      },
      onError: () => toastUnexpectedError("Failed to stop score engine"),
    });
  };
</script>

{#if confirmStopEngineId === instanceId}
  <wa-button
    size="s"
    appearance="plain"
    variant="neutral"
    onclick={() => handleCancel()}
  >
    Cancel
  </wa-button>
{/if}

<wa-button
  size="s"
  appearance="outlined"
  variant={confirmStopEngineId === instanceId ? "danger" : "neutral"}
  loading={stopScoreEngine.isPending}
  onclick={() => handleStop()}
>
  {#if confirmStopEngineId === instanceId}
    Confirm
  {:else}
    Stop
  {/if}

  <wa-icon name="stop" slot="start"></wa-icon>
</wa-button>
