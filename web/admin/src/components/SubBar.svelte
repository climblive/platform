<script lang="ts">
  interface Props {
    percentage: number;
    opacity: number;
  }

  const { percentage, opacity }: Props = $props();
</script>

<div
  class="sub-bar"
  style:--target-height="{percentage}%"
  style:--fill-weight="{opacity * 100}%"
></div>

<style>
  .sub-bar {
    width: 100%;
    background: color-mix(
      in oklab,
      var(--color) var(--fill-weight),
      var(--bar-fill-color)
    );
    animation: grow 0.6s ease-out forwards;
    height: 0;
    position: relative;
    flex: 0 0 auto;
  }

  @keyframes grow {
    to {
      height: var(--target-height);
    }
  }
</style>
