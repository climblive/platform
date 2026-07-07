<script lang="ts">
  type Props = {
    onClick: (e: MouseEvent) => void;
    iconName: string;
    label: string;
    points?: number;
    bonusPoints?: number;
    active: boolean;
  };

  const {
    onClick,
    iconName,
    label,
    points = 0,
    bonusPoints,
    active,
  }: Props = $props();

  const pointsLabel = $derived.by(() => {
    if (bonusPoints) {
      return `${points}p + ${bonusPoints}p`;
    }

    return `${points}p`;
  });
</script>

<div data-active={active}>
  <wa-button
    size="s"
    appearance="outlined"
    onclick={onClick}
    pill
    variant="neutral"
  >
    <wa-icon slot="start" name={iconName}></wa-icon>
    {label}
  </wa-button>
  {#if pointsLabel !== undefined}
    <span>{pointsLabel}</span>
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

  div[data-active="true"] {
    & span {
      color: var(--wa-color-success-fill-loud);
    }
  }

  span {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
    margin-inline-start: var(--wa-space-s);
  }
</style>
