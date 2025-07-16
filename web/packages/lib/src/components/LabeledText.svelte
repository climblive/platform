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
    font-size: var(--wa-font-size-x-small);
    line-height: var(--wa-line-height-normal);
  }

  p {
    margin: 0;
  }
</style>
