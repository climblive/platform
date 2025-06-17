<script lang="ts">
  import type { HTMLAttributes } from "svelte/elements";
  import { type WithElementRef } from "../../../utils";

  let {
    ref = $bindable(null),
    children,
    ...restProps
  }: WithElementRef<HTMLAttributes<HTMLTableSectionElement>> = $props();
</script>

<thead bind:this={ref} {...restProps}>
  {@render children?.()}
</thead>

<style>
  thead {
    height: 2rem;

    background-color: var(--sl-color-primary-600);
    color: white;
  }
</style>
