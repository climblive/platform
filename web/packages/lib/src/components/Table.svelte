<script lang="ts">
  import type { Snippet } from "svelte";

  type Props = {
    columns: string[];
    children: Snippet;
  };

  const { columns, children }: Props = $props();
</script>

<div>
  {#each columns as column}
    <div class="header">{column}</div>
  {/each}
  {@render children()}
</div>

<style>
  div {
    display: grid;
    grid-template-columns: repeat(3, minmax(100px, 1fr));
  }

  .header {
    border-bottom: 1px solid #ccc;
  }
</style>
