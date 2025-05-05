<script lang="ts">
  import type { Snippet } from "svelte";

  type Props = {
    children: Snippet;
  };

  const { children }: Props = $props();
</script>

<tr>
  {@render children()}
</tr>

<style>
  tr {
    height: 2rem;
    cursor: pointer;
  }

  tr:nth-child(even) {
    background-color: var(--sl-color-neutral-50);
  }

  tr:hover {
    background-color: var(--sl-color-neutral-100);
  }
</style>
