<script lang="ts">
  import type { HTMLTdAttributes } from "svelte/elements";
  import { type WithElementRef } from "../../../utils";

  let {
    ref = $bindable(null),
    children,
    ...restProps
  }: WithElementRef<HTMLTdAttributes> = $props();
</script>

<td bind:this={ref} {...restProps}>
  {@render children?.()}
</td>

<style>
  td:first-of-type {
    padding-left: var(--sl-spacing-medium);
  }

  td {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>
