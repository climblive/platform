<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import type WaCheckbox from "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import { ordinalSuperscript } from "@climblive/lib/utils";

  type Props = {
    onChange: (checked: boolean) => void;
    label: string;
    sublabel?: string;
    attempts: number;
    points?: number;
    bonusPoints?: number;
    checked?: boolean;
  };

  const {
    onChange,
    label,
    sublabel,
    attempts,
    points,
    bonusPoints,
    checked,
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

  const handleChange = (e: InputEvent) => {
    const checked = (e.target as WaCheckbox).checked;

    onChange(checked);
  };
</script>

<div class="container" data-active={checked}>
  <div class="top">
    <wa-checkbox onchange={handleChange} {checked} size="m">
      {label}
      {#if sublabel !== undefined}
        <span class="sublabel">
          {sublabel}
        </span>
      {/if}
    </wa-checkbox>
  </div>

  <div class="subtext">
    <span>
      {#if checked}
        {attempts}{ordinalSuperscript(attempts)}
        attempt
        <wa-icon name="lock"></wa-icon>
      {:else if attempts === 0}
        -
      {:else if attempts === 1}
        -
      {:else}
        -
      {/if}
    </span>
    {#if pointsLabel !== undefined}
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

  .top {
    display: flex;
    gap: var(--wa-space-s);
    width: 100%;
  }

  .sublabel {
    color: var(--wa-color-text-quiet);
  }
</style>
