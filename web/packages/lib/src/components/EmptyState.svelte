<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    title: string;
    description: string;
    actions?: Snippet;
  }

  const { title, description, actions }: Props = $props();
</script>

<div class="empty-state">
  <div class="content">
    <h3>{title}</h3>
    <p>{description}</p>
    {#if actions}
      <div class="actions">
        {@render actions()}
      </div>
    {/if}
  </div>
</div>

<style>
  .empty-state {
    width: 100%;
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-neutral-border-quiet);
    border-radius: var(--wa-border-radius-m);
    padding: var(--wa-space-2xl);
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .content {
    text-align: center;
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
    max-width: 32rem;
  }

  h3 {
    margin: 0;
    font-size: var(--wa-font-size-l);
    font-weight: var(--wa-font-weight-bold);
    color: var(--wa-color-text-default);
  }

  p {
    margin: 0;
    color: var(--wa-color-text-quiet);
    font-size: var(--wa-font-size-m);
  }

  .actions {
    margin-block-start: var(--wa-space-m);
    display: flex;
    justify-content: center;
    gap: var(--wa-space-xs);
  }
</style>
