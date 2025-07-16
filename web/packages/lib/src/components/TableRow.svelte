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
    background-color: var(--wa-color-gray-50);
  }

  tr:hover {
    background-color: var(--wa-color-gray-100);
  }
</style>
