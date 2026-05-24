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

<section data-disabled={disabled} data-checked={checked}>
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
  <span class="attempts"
    >{attempts === 1 ? "1 attempt" : `${attempts} attempts`}</span
  >
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
    align-items: center;

    &:has(wa-checkbox:disabled) {
      opacity: 0.5;
      cursor: not-allowed;
    }

    &[data-checked="true"] {
      border-color: var(--wa-color-green-50);
    }
  }

  .points {
    text-align: right;
  }

  .attempts {
    margin-block-start: var(--wa-space-xs);
    font-size: var(--wa-font-size-xs);
  }
</style>
