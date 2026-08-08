<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import type WaCheckbox from "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";

  type Props = {
    onChange: (checked: boolean, flash: boolean) => void;
    label: string;
    subLabel?: string;
    attempts: number;
    points?: number;
    bonusPoints?: number;
    checked?: boolean;
    indeterminate?: boolean;
    flashToggle?: boolean;
    showPoints?: boolean;
    showAttempts?: boolean;
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
    flashToggle,
    showPoints,
    showAttempts,
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

  const handleChange = (e: InputEvent, flash: boolean) => {
    const checked = (e.target as WaCheckbox).checked;

    onChange(checked, flash);
  };
</script>

<div class="container" data-active={checked} data-disabled={indeterminate}>
  <div class="top">
    <wa-checkbox
      onchange={(e: InputEvent) => handleChange(e, false)}
      checked={checked && !indeterminate && attempts > 1}
      {indeterminate}
      disabled={indeterminate}
      size="m"
    >
      {label}
    </wa-checkbox>

    {#if flashToggle}
      /

      <wa-checkbox
        onchange={(e: InputEvent) => handleChange(e, true)}
        checked={checked && !indeterminate && attempts === 1}
        {indeterminate}
        disabled={indeterminate}
        size="m"
      >
        Flash
        <wa-icon name="bolt"></wa-icon>
      </wa-checkbox>
    {/if}

    <span class="sub">{subLabel}</span>
  </div>

  <div class="subtext">
    {#if showAttempts}
      <span>
        {attempts}
        {attempts === 1 ? "attempt" : "attempts"}
        {#if checked}<wa-icon name="lock"></wa-icon>
        {/if}
      </span>
    {/if}
    {#if showPoints && pointsLabel !== undefined}
      <span>{pointsLabel}</span>
    {/if}
  </div>
</div>

<style>
  .container {
    width: 100%;
    display: flex;
    flex-direction: column;
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
    align-items: center;
  }

  .subtext {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
  }

  .subtext:not(:has(*)) {
    display: none;
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

  .sub {
    margin-left: auto;
    color: var(--wa-color-text-quiet);
  }

  .top {
    display: flex;
    gap: var(--wa-space-s);
    width: 100%;
  }
</style>
