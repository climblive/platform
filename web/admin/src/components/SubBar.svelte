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
    background: color-mix(
      in oklab,
      var(--color) var(--fill-weight),
      var(--bar-fill-color) 0%
    );
    animation: grow 0.6s ease-out forwards;
    height: 0;
  }

  @keyframes grow {
    to {
      height: var(--target-height);
    }
  }
</style>
