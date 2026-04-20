<script lang="ts">
  interface Props {
    percentage: number;
    color: string;
    opacity: number;
  }

  const { percentage, color, opacity }: Props = $props();
</script>

<div class="sub-bar-slot">
  <div class="sub-bar-bg" style:background={color}></div>
  <div
    class="sub-bar"
    style:--target-height="{percentage}%"
    style:background={color}
    style:opacity
  ></div>
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
