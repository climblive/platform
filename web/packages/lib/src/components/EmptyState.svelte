<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    title: string;
    description: string;
    actions?: Snippet;
  }

  const { title, description, actions }: Props = $props();
</script>

<section>
  <h3>{title}</h3>
  <p>{description}</p>

  {#if actions}
    {@render actions()}
  {/if}
</section>

<style>
  section {
    width: 100%;
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-neutral-border-quiet);
    border-radius: var(--wa-border-radius-m);
    padding: var(--wa-space-xl);
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
    justify-content: center;
    align-items: center;
    text-align: center;
  }

  h3 {
    margin: 0;
  }

  p {
    margin-block-start: 0;
    color: var(--wa-color-text-quiet);
    font-size: var(--wa-font-size-s);
  }
</style>
