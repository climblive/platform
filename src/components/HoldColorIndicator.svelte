<script lang="ts">
  export let primary: string;
  export let secondary: string | undefined = undefined;
</script>

<svg viewBox="0 0 100 100">
  <circle cx="50" cy="50" r="50" fill={primary} />
  {#if secondary}
    <path d="M0,50 a1,1 0 0,0 100,0" fill={secondary} />
  {/if}
  <circle cx="50" cy="50" r="45" stroke-width="10" fill="none" />
</svg>

<style>
  svg {
    height: 100%;
    width: 100%;
    transform: rotate(-45deg);

    & circle {
      stroke: var(--sl-color-primary-700);
    }
  }
</style>
