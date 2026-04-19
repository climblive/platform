<script lang="ts">
  interface Props {
    primary: string | undefined;
    secondary?: string | undefined;
  }

  let { primary, secondary }: Props = $props();
</script>

<div class="indicator">
  <div class="fill" class:checkerboard={primary === undefined}>
    {#if primary || secondary}
      <svg viewBox="0 0 100 100">
        {#if primary}
          <circle cx="50" cy="50" r="50" fill={primary} />
        {/if}
        {#if secondary}
          <path d="M0,50 a1,1 0 0,0 100,0" fill={secondary} />
        {/if}
      </svg>
    {/if}
  </div>
</div>

<style>
  .indicator {
    height: var(--height, 100%);
    width: var(--width, 100%);
    border-radius: var(--wa-border-radius-pill);
    box-sizing: border-box;
    overflow: hidden;
    position: relative;
  }

  .indicator::after {
    content: "";
    position: absolute;
    inset: 0;
    border-radius: var(--wa-border-radius-pill);
    outline: var(--wa-form-control-border-width)
      var(--wa-form-control-border-style) var(--wa-form-control-border-color);
    outline-offset: calc(
      -1 * var(--wa-form-control-border-width) - var(--wa-space-3xs)
    );
    pointer-events: none;
  }

  .fill {
    height: 100%;
    width: 100%;
    border-radius: var(--wa-border-radius-circle);
    overflow: hidden;
  }

  .checkerboard {
    background-image:
      linear-gradient(
        45deg,
        var(--wa-color-neutral-fill-normal) 25%,
        transparent 25%
      ),
      linear-gradient(
        45deg,
        transparent 75%,
        var(--wa-color-neutral-fill-normal) 75%
      ),
      linear-gradient(
        45deg,
        transparent 75%,
        var(--wa-color-neutral-fill-normal) 75%
      ),
      linear-gradient(
        45deg,
        var(--wa-color-neutral-fill-normal) 25%,
        transparent 25%
      );
    background-size: 8px 8px;
    background-position:
      0 0,
      0 0,
      -4px -4px,
      4px 4px;
  }

  svg {
    display: block;
    height: 100%;
    width: 100%;
    transform: rotate(-45deg);
  }
</style>
