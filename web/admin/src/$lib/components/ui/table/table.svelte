<script lang="ts">
  import type { HTMLTableAttributes } from "svelte/elements";
  import { type WithElementRef } from "../../../utils";

  let {
    ref = $bindable(null),
    children,
    ...restProps
  }: WithElementRef<HTMLTableAttributes> = $props();
</script>

<table bind:this={ref} {...restProps}>
  {@render children?.()}
</table>

<style>
  table {
    margin-top: 1rem;
    width: 100%;
    table-layout: fixed;

    border: none;
    border-collapse: separate;
    overflow: hidden;
    border-spacing: 0;
    border-radius: var(--sl-border-radius-medium);
  }
</style>
