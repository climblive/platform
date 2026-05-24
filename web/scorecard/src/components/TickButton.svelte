<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";

  type Props = {
    onClick: (e: MouseEvent) => void;
    onAdjustAttempts: (e: MouseEvent, delta: number) => void;
    label: string;
    points?: number;
    attempts: number;
    checked: boolean;
    indeterminate?: boolean;
    disabled?: boolean;
  };

  const {
    onClick,
    onAdjustAttempts,
    label,
    points,
    attempts,
    checked,
    indeterminate = false,
    disabled = false,
  }: Props = $props();
</script>

<section data-checked={checked}>
  <wa-checkbox
    class:checked-state={checked}
    size="s"
    {checked}
    {indeterminate}
    {disabled}
    aria-label={label}
    onclick={(event: Event) => {
      if (disabled) {
        event.preventDefault();
        return;
      }

      event.preventDefault();
      onClick(event as MouseEvent);
    }}
  >
    {label}
  </wa-checkbox>

  {#if points !== undefined}
    <span class="points">{points}p</span>
  {/if}
  <div class="attempt-controls">
    <span class="attempts"
      >{attempts === 1 ? "1 attempt" : `${attempts} attempts`}</span
    >
    <wa-button
      size="xs"
      appearance="outlined"
      aria-label={`Decrease ${label} attempts`}
      onclick={(event: MouseEvent) => onAdjustAttempts(event, -1)}
    >
      <wa-icon name="minus" label={`Decrease ${label} attempts`}></wa-icon>
    </wa-button>
    <wa-button
      size="xs"
      appearance="outlined"
      aria-label={`Increase ${label} attempts`}
      onclick={(event: MouseEvent) => onAdjustAttempts(event, 1)}
    >
      <wa-icon name="plus" label={`Increase ${label} attempts`}></wa-icon>
    </wa-button>
  </div>
</section>

<style>
  section {
    display: grid;
    grid-template-columns: max-content 1fr;

    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-neutral-border-normal);
    border-radius: var(--wa-border-radius-m);
    background-color: var(--wa-color-surface-raised);
    padding: var(--wa-space-s);
    gap: var(--wa-space-s);
    align-items: center;

    &[data-checked="true"] {
      border-color: var(--wa-color-green-50);
    }
  }

  .points {
    text-align: right;
  }

  .attempts {
    font-size: var(--wa-font-size-xs);
  }

  .attempt-controls {
    display: flex;
    align-items: center;
    gap: var(--wa-space-3xs);
    margin-block-start: var(--wa-space-xs);
    margin-inline-start: 1.5rem;
  }
</style>
