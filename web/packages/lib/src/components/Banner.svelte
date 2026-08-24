<script lang="ts">
  import type { Snippet } from "svelte";

  type Props = {
    visible: boolean;
    children: Snippet;
  };

  const { visible, children }: Props = $props();
</script>

<header data-visible={visible}>
  {@render children()}
</header>

<style>
  header {
    background-color: var(--wa-color-danger-fill-loud);
    width: 100%;
    height: 2rem;
    display: flex;
    justify-content: center;
    align-items: center;
    gap: var(--wa-space-xs);
    color: var(--wa-color-danger-on-loud);
    font-weight: var(--wa-font-weight-semibold);
    font-size: var(--wa-font-size-s);

    &[data-visible="false"] {
      display: none;
    }
  }
</style>
