<script lang="ts">
  interface Props {
    primary?: string;
    secondary?: string | undefined;
  }

  let { primary, secondary }: Props = $props();
</script>

<svg viewBox="0 0 100 100">
  <circle cx="50" cy="50" r="50" fill={primary} />
  {#if secondary}
    <path d="M0,50 a1,1 0 0,0 100,0" fill={secondary} />
  {/if}
  <circle
    cx="50"
    cy="50"
    r="45"
    fill="none"
    stroke="var(--wa-color-text-normal)"
    stroke-width="10"
  />
</svg>

<style>
  svg {
    height: var(--height, 100%);
    width: var(--width, 100%);
    transform: rotate(-45deg);
  }
</style>
