<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";

  type Props = {
    onClick: (e: MouseEvent) => void;
    label: string;
    points?: number;
    attempts: number;
    checked: boolean;
    indeterminate?: boolean;
    disabled?: boolean;
  };

  const {
    onClick,
    label,
    points,
    attempts,
    checked,
    indeterminate = false,
    disabled = false,
  }: Props = $props();
</script>

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
  <span class="content">
    <span class="label">{label}</span>
    {#if points !== undefined}
      <span class="points">{points}p</span>
    {/if}
    <span class="attempts"
      >{attempts === 1 ? "1 attempt" : `${attempts} attempts`}</span
    >
  </span>
</wa-checkbox>

<style>
  wa-checkbox {
    width: 100%;
    display: block;

    &::part(base) {
      border: var(--wa-border-width-s) var(--wa-border-style)
        var(--wa-color-neutral-border-normal);
      border-radius: var(--wa-border-radius-m);
      background-color: var(--wa-color-surface-raised);
      padding: var(--wa-space-s);
      min-height: 3rem;
      box-sizing: border-box;
      align-items: center;
    }

    &::part(label) {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;
      min-width: 0;
    }

    &.checked-state::part(base) {
      border-color: var(--wa-color-green-50);
    }
  }

  .content {
    width: 100%;
    display: grid;
    grid-template-columns: max-content max-content;
    grid-template-areas:
      "label points"
      "attempts points";
    align-items: center;
    justify-content: center;
    column-gap: var(--wa-space-s);
    row-gap: var(--wa-space-3xs);
    font-size: var(--wa-font-size-s);
  }

  .label {
    grid-area: label;
    font-weight: var(--wa-font-weight-medium);
    text-align: center;
  }

  .attempts {
    grid-area: attempts;
    justify-self: start;
  }

  .points {
    grid-area: points;
    align-self: center;
  }

  .attempts,
  .points {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
    white-space: nowrap;
  }
</style>
