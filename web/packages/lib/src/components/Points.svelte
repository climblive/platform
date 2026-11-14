<script lang="ts">
  import { animatePoints } from "../utils/animatePoints";
  import { pointsNumberFormat, roundPoints } from "../utils/points";

  interface Props {
    value: number;
    compactThreshold?: number;
  }

  const { value, compactThreshold = Number.MAX_SAFE_INTEGER }: Props = $props();

  const formatter = $derived(pointsNumberFormat(value, compactThreshold));
  const formattedValue = $derived(formatter.format(roundPoints(value)) + "p");
  const animatedFormatter = $derived.by(() => {
    const fractionDigits =
      formatter
        .formatToParts(roundPoints(value))
        .find((part) => part.type === "fraction")?.value.length ?? 0;

    return new Intl.NumberFormat(undefined, {
      ...formatter.resolvedOptions(),
      minimumFractionDigits: fractionDigits,
      maximumFractionDigits: fractionDigits,
    });
  });

  let element: HTMLSpanElement;
  let animated = $state<number>();
  let previous = 0;

  $effect(() => {
    const target = value;
    const from = previous;
    previous = target;

    const stop = animatePoints(element, from, target, (value) => {
      animated = value;
    });

    return () => {
      const current = animated;
      stop();
      previous = current ?? target;
    };
  });
</script>

<span class="points" bind:this={element}>
  <span class="value" class:animating={animated !== undefined}>
    {formattedValue}
  </span>
  {#if animated !== undefined}
    <span class="animation" aria-hidden="true">
      {animatedFormatter.format(roundPoints(animated))}p
    </span>
  {/if}
</span>

<style>
  .points {
    display: inline-grid;
    font-variant-numeric: tabular-nums;
  }

  .value,
  .animation {
    grid-area: 1 / 1;
  }

  .animation {
    display: none;
  }

  @media (prefers-reduced-motion: no-preference) {
    .value.animating {
      opacity: 0;
    }

    .animation {
      display: inline;
    }
  }
</style>
