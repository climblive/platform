<script lang="ts">
  import type { Snippet } from "svelte";

  type Props = {
    children: Snippet;
  };

  const { children }: Props = $props();
</script>

<div>
  {@render children()}
</div>

<style>
  div {
    display: grid;
    grid-column: 1 / -1;
    grid-template-columns: subgrid;
    grid-template-rows: 1fr;

    height: 2rem;
    align-items: center;
  }

  div:nth-child(odd) {
    background-color: var(--sl-color-neutral-100);
  }

  div {
    padding-inline: 2rem;
  }
</style>
