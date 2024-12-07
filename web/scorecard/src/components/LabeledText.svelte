<script lang="ts">
  export let label: string;
  export let text: string;
</script>

<section>
  <p>
    <strong>{label}</strong>
    {text}
  </p>
</section>

<style>
  strong {
    display: block;
    font-size: var(--sl-font-size-small);
  }
</style>
