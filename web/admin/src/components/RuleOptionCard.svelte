<script lang="ts">
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import "@awesome.me/webawesome/dist/components/radio/radio.js";
  import type { Snippet } from "svelte";

  interface Props {
    title: string;
    description: string;
    disabled?: boolean;
    tag?: string;
    children: Snippet;
    footer?: Snippet;
  }

  const {
    title,
    description,
    disabled = false,
    tag,
    children,
    footer,
  }: Props = $props();
</script>

<div class="card" data-disabled={disabled}>
  <div class="header">
    {@render children()}
    <h3>{title}</h3>
    {#if tag}
      <wa-badge pill variant="neutral">{tag}</wa-badge>
    {/if}
  </div>

  <p class="description">{description}</p>

  {#if footer}
    {@render footer()}
  {/if}
</div>

<style>
  .card {
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-neutral-border-normal);
    border-radius: var(--wa-border-radius-m);
    padding: var(--wa-space-m);
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);

    &[data-disabled="true"] {
      opacity: 0.5;
      pointer-events: none;
    }
  }

  .header {
    display: flex;
    align-items: center;
    gap: var(--wa-space-s);

    & h3 {
      margin: 0;
      font-size: var(--wa-font-size-m);
    }

    wa-badge {
      margin-inline-start: auto;
      font-size: var(--wa-font-size-2xs);
    }
  }

  .description {
    margin: 0;
    font-size: var(--wa-font-size-s);
    color: var(--wa-color-text-quiet);
  }
</style>
