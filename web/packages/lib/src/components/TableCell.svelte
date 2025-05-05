<script lang="ts">
  import type { Snippet } from "svelte";

  type Props = {
    children: Snippet;
    align?: "left" | "right";
  };

  const { children, align = "left" }: Props = $props();
</script>

<td data-align={align}>
  {@render children()}
</td>

<style>
  td:first-of-type {
    padding-left: var(--sl-spacing-medium);
  }

  td[data-align="right"] {
    text-align: right;
  }
</style>
