<script lang="ts">
  import type { HTMLAttributes } from "svelte/elements";
  import { type WithElementRef } from "../../../utils";

  let {
    ref = $bindable(null),
    children,
    ...restProps
  }: WithElementRef<HTMLAttributes<HTMLTableSectionElement>> = $props();
</script>

<tbody bind:this={ref} {...restProps}>
  {@render children?.()}
</tbody>
