<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { onMount, type Snippet } from "svelte";

  interface Props {
    readyTitle: string;
    narrowMargins?: boolean;
    preparePrint?: () => void;
    children: Snippet;
  }

  const {
    readyTitle,
    narrowMargins = false,
    preparePrint,
    children,
  }: Props = $props();

  let printingFinished = $state(false);

  onMount(() => {
    const timer = window.setTimeout(() => window.print());
    return () => window.clearTimeout(timer);
  });
</script>

<svelte:window
  onbeforeprint={preparePrint}
  onafterprint={() => (printingFinished = true)}
/>

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
    <wa-button variant="brand" onclick={() => window.print()} size="s">
      <wa-icon slot="start" name="print"></wa-icon>
      Print again
    </wa-button>
  {:else}
    <h1>{readyTitle}</h1>
    <p>If the print dialog doesn't open automatically, click Print below.</p>
    <wa-button variant="brand" onclick={() => window.print()} size="s">
      <wa-icon slot="start" name="print"></wa-icon>
      Print
    </wa-button>
  {/if}
</section>

<div class="print-content" class:narrow={narrowMargins}>
  {@render children()}
</div>

<style>
  .print-content {
    position: absolute;
    inset-inline-start: -10000px;
    visibility: hidden;
    page: normal-print;

    &.narrow {
      page: narrow-print;
    }
  }

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

    .print-content {
      position: static;
      visibility: visible;
    }

    :global(html),
    :global(body) {
      background: white;
      color: black;
      overflow: visible;
      transform: none;
      will-change: auto;
    }
  }

  @page normal-print {
    size: a4 portrait;
    margin: 20mm;
  }

  @page narrow-print {
    size: a4 portrait;
    margin: 10mm;
  }
</style>
