<script lang="ts">
  interface Props {
    percentage: number;
    color: string;
    opacity: number;
    visible: boolean;
  }

  const { percentage, color, opacity, visible }: Props = $props();

  const minHeight = "2px";
</script>

<div class="sub-bar-slot">
  {#if visible}
    <div class="sub-bar-bg" style:background={color}></div>
    <div
      class="sub-bar"
      style:--target-height={percentage > 0 ? `${percentage}%` : minHeight}
      style:background={color}
      style:opacity
    ></div>
  {/if}
</div>

<style>
  .sub-bar-slot {
    flex: 1;
    position: relative;
    height: 100%;
    display: flex;
    align-items: flex-end;
  }

  .sub-bar-bg {
    position: absolute;
    inset: 0;
    border-radius: var(--wa-border-radius-s) var(--wa-border-radius-s) 0 0;
    opacity: 0.15;
  }

  .sub-bar {
    width: 100%;
    border-radius: var(--wa-border-radius-s);
    animation: grow 0.6s ease-out forwards;
    height: 0;
    position: relative;
  }

  @keyframes grow {
    to {
      height: var(--target-height);
    }
  }
</style>
