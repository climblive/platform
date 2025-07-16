<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    order: number;
    children?: Snippet;
  }

  let { order, children }: Props = $props();
</script>

<div style="--order: {order}">
  {@render children?.()}
</div>

<style>
  div {
    position: absolute;
    left: 0;
    right: 0;
    transition: transform 0.3s ease-out;
    transform: translateY(
      calc(
        (var(--order) - var(--page-index) * var(--page-size)) *
          (2.25rem + var(--wa-spacing-x-small))
      )
    );
  }
</style>
