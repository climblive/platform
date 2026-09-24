<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  type Props = {
    onClick: () => void;
    iconName: string;
    label: string;
    points?: number;
    bonusPoints?: number;
    disabled?: boolean;
    reached?: boolean;
  };

  const {
    onClick,
    iconName,
    label,
    points,
    bonusPoints,
    disabled,
    reached = false,
  }: Props = $props();

  const pointsLabel = $derived.by(() => {
    if (points === undefined) {
      return "-";
    }

    if (bonusPoints) {
      return `${points}p + ${bonusPoints}p`;
    }

    return `${points}p`;
  });
</script>

<div>
  <wa-button
    size="s"
    appearance={reached ? "filled-outlined" : "outlined"}
    onclick={onClick}
    {disabled}
    pill
    variant={reached ? "success" : "neutral"}
  >
    {#if reached}
      <wa-icon slot="start" name={iconName}></wa-icon>
    {/if}
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

  span {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
    margin-inline-start: var(--wa-space-s);
  }
</style>
