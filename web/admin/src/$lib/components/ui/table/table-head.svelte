<script lang="ts">
  import type { HTMLThAttributes } from "svelte/elements";
  import { type WithElementRef } from "../../../utils";

  let {
    ref = $bindable(null),
    children,
    ...restProps
  }: WithElementRef<HTMLThAttributes> = $props();
</script>

<th bind:this={ref} {...restProps}>
  {@render children?.()}
</th>

<style>
  th {
    font-weight: var(--sl-font-weight-normal);
    text-align: left;
  }

  th:first-of-type {
    padding-left: var(--sl-spacing-medium);
  }
</style>
