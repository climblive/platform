<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";

  type Props = {
    onClick: (e: MouseEvent) => void;
    label: string;
    attempts?: number;
    points?: number;
    bonusPoints?: number;
    checked: boolean;
    indeterminate?: boolean;
  };

  const {
    onClick,
    label,
    attempts,
    points = 0,
    bonusPoints,
    checked,
    indeterminate = false,
  }: Props = $props();

  const pointsLabel = $derived.by(() => {
    if (bonusPoints) {
      return `${points}p + ${bonusPoints}p`;
    }

    return `${points}p`;
  });
</script>

<div class="container" data-active={checked}>
  <wa-checkbox
    onclick={onClick}
    {checked}
    {indeterminate}
    disabled={indeterminate}
  >
    {label}
  </wa-checkbox>
  <div>
    {#if attempts !== undefined}
      <span>{attempts} attempts;</span>
    {/if}
    {#if pointsLabel !== undefined}
      <span>{pointsLabel}</span>
    {/if}
  </div>
</div>

<style>
  wa-checkbox {
    width: 100%;
  }

  .container {
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: start;
    gap: var(--wa-space-2xs);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-surface-border);
    padding: var(--wa-space-s);
    border-radius: var(--wa-border-radius-m);
  }

  .container[data-active="true"] {
    & span {
      color: var(--wa-color-success-fill-loud);
    }
    border-color: var(--wa-color-success-border-loud);
  }

  span {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
  }
</style>
