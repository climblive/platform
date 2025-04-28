<script lang="ts">
  import type { Snippet } from "svelte";

  type Props = {
    columns: string[];
    children: Snippet;
  };

  const { columns, children }: Props = $props();
</script>

<div class="table" style={`--columns: ${columns.length}`}>
  <div class="header">
    {#each columns as column}
      <div>{column}</div>
    {/each}
  </div>
  {@render children()}
</div>

<style>
  .table {
    margin-top: 1rem;
    width: 100%;
    display: grid;
    grid-template-columns: repeat(var(--columns), minmax(min-content, 1fr));
  }

  .header {
    display: grid;
    grid-column: 1 / -1;
    grid-template-columns: subgrid;
    grid-template-rows: 1fr;

    height: 2rem;
    align-items: center;
    padding-inline: 2rem;

    font-weight: var(--sl-font-weight-bold);
    background-color: var(--sl-color-primary-600);
    border-top-left-radius: var(--sl-border-radius-medium);
    border-top-right-radius: var(--sl-border-radius-medium);
    color: white;
  }
</style>
