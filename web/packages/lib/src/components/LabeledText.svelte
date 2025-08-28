<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    label: string;
    children?: Snippet;
  }

  let { label, children }: Props = $props();
</script>

<section>
  <p>
    <strong>{label}</strong>
    {@render children?.()}
  </p>
</section>

<style>
  strong {
    display: block;
    font-size: var(--wa-font-size-xs);
    line-height: var(--wa-line-height-normal);
    margin-block-end: var(--wa-space-2xs);
    color: var(--wa-color-text-quiet);
  }

  p {
    margin: 0;
  }
</style>
