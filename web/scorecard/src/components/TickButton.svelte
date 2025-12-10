<script lang="ts">
  type Props = {
    onClick: (e: MouseEvent) => void;
    iconName: string;
    label: string;
    className?: string;
    points?: number;
    flash?: boolean;
  };

  const { onClick, iconName, label, className, points, flash }: Props =
    $props();
</script>

<div data-flash={flash}>
  <wa-button
    size="small"
    appearance="outlined"
    class={className}
    onclick={onClick}
    pill
  >
    <wa-icon slot="start" name={iconName}></wa-icon>
    {label}
  </wa-button>
  {#if points}
    <span>+{points} pts</span>
  {/if}
</div>

<style>
  wa-button {
    width: 100%;
  }

  div {
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: start;
    gap: var(--wa-space-2xs);

    &[data-flash="true"] wa-button {
      --wa-color-fill-quiet: var(--wa-color-yellow-10);
      --wa-color-border-loud: var(--wa-color-yellow-40);
      --wa-color-border-normal: var(--wa-color-yellow-30);
      --wa-color-border-quiet: var(--wa-color-yellow-20);
      --wa-color-on-normal: var(--wa-color-yellow-70);
      --wa-color-on-quiet: var(--wa-color-yellow-60);
    }
  }

  span {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
    margin-inline-start: var(--wa-space-s);
  }
</style>
