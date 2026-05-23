<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";

  type Props = {
    onClick: (e: MouseEvent) => void;
    label: string;
    points?: number;
    active: boolean;
  };

  const { onClick, label, points, active }: Props = $props();
</script>

<wa-checkbox
  size="s"
  checked={active}
  data-active={active}
  aria-label={label}
  onclick={(event: Event) => {
    event.preventDefault();
    onClick(event as MouseEvent);
  }}
>
  <span class="content">
    <span class="label">{label}</span>
    {#if points !== undefined}
      <span class="points">{points}p</span>
    {/if}
  </span>
</wa-checkbox>

<style>
  wa-checkbox {
    width: 100%;
    display: block;
    --checked-icon-color: var(--wa-color-success-fill-loud);

    &::part(base) {
      border: var(--wa-border-width-s) var(--wa-border-style)
        var(--wa-color-neutral-border-normal);
      border-radius: var(--wa-border-radius-m);
      background-color: var(--wa-color-surface-raised);
      padding: var(--wa-space-s);
      min-height: 3rem;
      box-sizing: border-box;
      transition:
        border-color var(--wa-transition-fast),
        background-color var(--wa-transition-fast);
    }

    &::part(label) {
      flex: 1;
      min-width: 0;
    }

    &[data-active="true"]::part(base) {
      border-color: var(--wa-color-success-fill-loud);
      background-color: var(--wa-color-success-95);
    }
  }

  .content {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--wa-space-s);
    font-size: var(--wa-font-size-s);
  }

  .label {
    font-weight: var(--wa-font-weight-medium);
  }

  .points {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
    white-space: nowrap;
  }

  wa-checkbox[data-active="true"] .points {
    color: var(--wa-color-success-fill-loud);
  }
</style>
