<script lang="ts">
  export let label: string;
</script>

<section>
  <p>
    <strong>{label}</strong>
    <slot />
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
