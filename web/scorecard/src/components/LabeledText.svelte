<script lang="ts">
  interface Props {
    label: string;
    children?: import('svelte').Snippet;
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
    font-size: var(--sl-font-size-x-small);
    line-height: var(--sl-line-height-normal);
  }

  p {
    margin: 0;
  }
</style>
