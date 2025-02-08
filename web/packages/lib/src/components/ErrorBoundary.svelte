<script lang="ts">
  import "@shoelace-style/shoelace/dist/components/button/button.js";
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
      <sl-button size="small" variant="primary" onclick={reset}
        >Try again</sl-button
      >
    </main>
  {/snippet}
</svelte:boundary>

<style>
  pre {
    border: 1px solid var(--sl-color-primary-300);
    background-color: white;
    padding: var(--sl-spacing-small);
    overflow: scroll;
    border-radius: var(--sl-border-radius-small);
    font-size: var(--sl-font-size-small);
  }

  main {
    padding: var(--sl-spacing-small);
  }
</style>
