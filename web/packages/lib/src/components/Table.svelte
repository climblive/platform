<script lang="ts">
  import type { Snippet } from "svelte";

  type Props = {
    columns: string[];
    children: Snippet;
  };

  const { columns, children }: Props = $props();
</script>

<div style={`--columns: ${columns.length}`}>
  {#each columns as column}
    <div class="header">{column}</div>
  {/each}
  {@render children()}
</div>

<style>
  div {
    width: 100%;
    display: grid;
    grid-template-columns: repeat(var(--columns), minmax(100px, 1fr));
  }

  .header {
    border-bottom: 1px solid #ccc;
  }
</style>
