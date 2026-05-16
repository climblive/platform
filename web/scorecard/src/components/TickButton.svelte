<script lang="ts">
  type Props = {
    onClick: (e: MouseEvent) => void;
    iconName: string;
    label: string;
    points?: number;
    variant?: "danger" | "success";
  };

  const { onClick, iconName, label, points, variant }: Props = $props();
</script>

<div>
  <wa-button
    size="s"
    appearance={variant === "danger" ? "plain" : "outlined"}
    {variant}
    onclick={onClick}
    pill
  >
    <wa-icon slot="start" name={iconName}></wa-icon>
    {label}
  </wa-button>
  {#if points}
    <span data-variant={variant}>+{points} pts</span>
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

    &[data-variant="danger"] {
      color: var(--wa-color-brand-fill-loud);
    }

    &[data-variant="success"] {
      color: var(--wa-color-success-fill-loud);
    }
  }
</style>
