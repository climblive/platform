<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { onMount, tick } from "svelte";

  interface Props {
    readyTitle: string;
    preparePrint?: () => void | Promise<void>;
  }

  const { readyTitle, preparePrint }: Props = $props();

  let printingFinished = $state(false);
  let preparing = $state(false);
  let active = false;

  const handlePrint = async () => {
    if (preparing) {
      return;
    }

    preparing = true;

    try {
      await tick();
      await document.fonts.ready;

      if (!active) {
        return;
      }

      await preparePrint?.();

      if (active) {
        window.print();
      }
    } finally {
      preparing = false;
    }
  };

  onMount(() => {
    active = true;
    void handlePrint();

    return () => {
      active = false;
    };
  });
</script>

<svelte:window onafterprint={() => (printingFinished = true)} />

<section class="print-prompt">
  {#if printingFinished}
    <h1>All done?</h1>
    <p>
      Print again, or close this tab to return to managing your competition.
    </p>
    <wa-button
      variant="neutral"
      appearance="outlined"
      onclick={() => window.close()}
      size="s"
    >
      <wa-icon slot="start" name="close"></wa-icon>
      Close
    </wa-button>
    <wa-button
      variant="brand"
      onclick={handlePrint}
      disabled={preparing}
      size="s"
    >
      <wa-icon slot="start" name="print"></wa-icon>
      Print again
    </wa-button>
  {:else}
    <h1>{readyTitle}</h1>
    <p>If the print dialog doesn't open automatically, click Print below.</p>
    <wa-button
      variant="brand"
      onclick={handlePrint}
      disabled={preparing}
      size="s"
    >
      <wa-icon slot="start" name="print"></wa-icon>
      Print
    </wa-button>
  {/if}
</section>

<style>
  .print-prompt {
    text-align: center;

    & p {
      color: var(--wa-color-text-quiet);
    }

    & wa-button + wa-button {
      margin-inline-start: var(--wa-space-xs);
    }
  }

  @media print {
    .print-prompt {
      display: none;
    }
  }
</style>
