<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import type WaNumberInput from "@awesome.me/webawesome/dist/components/number-input/number-input.js";

  type Props = {
    onClick: (e: MouseEvent) => void;
    onSaveAttempts: (e: MouseEvent, attempts: number) => void;
    label: string;
    points?: number;
    attempts: number;
    minimumAttempts: number;
    checked: boolean;
    indeterminate?: boolean;
    disabled?: boolean;
  };

  const {
    onClick,
    onSaveAttempts,
    label,
    points,
    attempts,
    minimumAttempts,
    checked,
    indeterminate = false,
    disabled = false,
  }: Props = $props();

  let editingAttempts = $state(false);
  let attemptInput: WaNumberInput | undefined = $state();

  const attemptLabel = $derived.by(() => {
    if (minimumAttempts > 0) {
      if (attempts === 1) {
        return "Flash";
      }

      return `${attempts} attempts`;
    }

    if (attempts === 1) {
      return "1 failed attempt";
    }

    if (attempts > 1) {
      return `${attempts} failed attempts`;
    }

    return "0 attempts";
  });

  const handleEditAttempts = (event: MouseEvent) => {
    event.stopPropagation();
    editingAttempts = true;

    queueMicrotask(() => attemptInput?.focus());
  };

  const handleSaveAttempts = (event: MouseEvent) => {
    event.stopPropagation();

    const nextAttempts = Number(attemptInput?.value ?? attempts);

    editingAttempts = false;
    onSaveAttempts(event, nextAttempts);
  };
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
    {#if editingAttempts}
      <wa-number-input
        bind:this={attemptInput}
        size="xs"
        appearance="outlined"
        min={minimumAttempts}
        step="1"
        value={attempts}
        aria-label={`${label} attempts`}
      ></wa-number-input>
      <wa-button
        size="xs"
        appearance="outlined"
        aria-label={`Save ${label} attempts`}
        onclick={handleSaveAttempts}
      >
        Save
      </wa-button>
    {:else}
      <span class="attempts">{attemptLabel}</span>
      <wa-button
        size="xs"
        appearance="plain"
        aria-label={`Edit ${label} attempts`}
        onclick={handleEditAttempts}
      >
        <wa-icon name="pen" label={`Edit ${label} attempts`}></wa-icon>
      </wa-button>
    {/if}
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
    gap: var(--wa-space-2xs);
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
    margin-inline-end: var(--wa-space-s);
  }

  .attempt-controls {
    display: flex;
    align-items: center;
    gap: var(--wa-space-2xs);
  }
</style>
