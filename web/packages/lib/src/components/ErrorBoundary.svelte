<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type { Snippet } from "svelte";

  interface Props {
    children?: Snippet;
  }

  let { children }: Props = $props();

  const copyToClipboard = async (error: unknown) => {
    await navigator.clipboard.writeText(error as string);
  };
</script>

<svelte:boundary>
  {@render children?.()}

  {#snippet failed(error, reset)}
    <main>
      <h1>Oopsie!</h1>
      <pre onclick={() => copyToClipboard(error)}>{error}</pre>
      <wa-button size="small" variant="primary" onclick={reset}
        >Try again</wa-button
      >
    </main>
  {/snippet}
</svelte:boundary>

<style>
  pre {
    border: 1px solid var(--wa-color-primary-300);
    background-color: white;
    padding: var(--wa-space-s);
    overflow: scroll;
    border-radius: var(--wa-border-radius-s);
    font-size: var(--wa-font-size-s);
  }

  main {
    padding: var(--wa-space-s);
  }
</style>
