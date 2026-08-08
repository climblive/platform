<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";

  type Props = {
    onChange: (e: InputEvent) => void;
    label: string;
    subLabel?: string;
    attempts?: number;
    points?: number;
    bonusPoints?: number;
    checked?: boolean;
    indeterminate?: boolean;
  };

  const {
    onChange,
    label,
    subLabel,
    attempts,
    points,
    bonusPoints,
    checked,
    indeterminate = false,
  }: Props = $props();

  const pointsLabel = $derived.by(() => {
    if (points === undefined) {
      return undefined;
    }

    if (bonusPoints) {
      return `${points}p + ${bonusPoints}p`;
    }

    return `${points}p`;
  });
</script>

<div class="container" data-active={checked} data-disabled={indeterminate}>
  <wa-checkbox
    onchange={onChange}
    checked={checked && !indeterminate}
    {indeterminate}
    disabled={indeterminate}
  >
    {label}
    <span class="sub">{subLabel}</span>
  </wa-checkbox>
  <div class="subtext">
    {#if attempts !== undefined}
      <span>
        {attempts}
        {attempts === 1 ? "attempt" : "attempts"}
        {#if checked}<wa-icon name="lock"></wa-icon>
        {/if}
      </span>
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
    gap: var(--wa-space-xs);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-surface-border);
    padding: var(--wa-space-s);
    border-radius: var(--wa-border-radius-m);
  }

  .subtext {
    display: flex;
    width: 100%;
    justify-content: space-between;
  }

  .subtext {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
  }

  .container[data-active="true"] {
    & .subtext {
      color: var(--wa-color-success-fill-loud);
    }

    border-color: var(--wa-color-success-border-loud);
  }

  .container[data-disabled="true"] {
    & .subtext {
      color: var(--wa-color-text-quiet);
    }

    border-color: var(--wa-color-surface-border);
    opacity: 0.5;
  }

  wa-checkbox::part(label) {
    width: 100%;
    display: flex;
    justify-content: space-between;
  }

  .sub {
    margin-left: auto;
    color: var(--wa-color-text-quiet);
  }
</style>
