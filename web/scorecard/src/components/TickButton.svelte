<script lang="ts">
  type Props = {
    onClick: (e: MouseEvent) => void;
    iconName: string;
    label: string;
    points?: number;
    flash?: boolean;
  };

  const { onClick, iconName, label, points, flash }: Props = $props();
</script>

<div>
  <wa-button size="small" appearance="outlined" onclick={onClick} pill>
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
    gap: var(--wa-space-xs);
  }

  span {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
    margin-inline-start: var(--wa-space-s);
  }
</style>
