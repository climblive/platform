<script lang="ts">
  import type { HTMLAttributes } from "svelte/elements";
  import { type WithElementRef } from "../../../utils";

  let {
    ref = $bindable(null),
    children,
    ...restProps
  }: WithElementRef<HTMLAttributes<HTMLTableRowElement>> = $props();
</script>

<tr bind:this={ref} {...restProps}>
  {@render children?.()}
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
