<script lang="ts">
  interface Props {
    percentage: number;
    fillWeight: number;
  }

  const { percentage, fillWeight }: Props = $props();
</script>

<div
  style:--target-height="{percentage}%"
  style:--fill-weight="{fillWeight * 100}%"
></div>

<style>
  div {
    background: rgb(from var(--color) r g b / var(--fill-weight));
    animation: grow var(--wa-transition-slow) var(--wa-transition-easing)
      forwards;
    height: 0;
  }

  @keyframes grow {
    to {
      height: var(--target-height);
    }
  }
</style>
