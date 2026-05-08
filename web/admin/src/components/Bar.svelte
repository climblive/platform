<script lang="ts">
  import "@awesome.me/webawesome/dist/components/popover/popover.js";
  import type { Problem } from "@climblive/lib/models";
  import type { ProblemStats } from "./ProblemsChart.svelte";
  import SubBar from "./SubBar.svelte";

  type Props = {
    problem: Problem;
    stats: ProblemStats;
    maxCount: number;
  };

  const { problem, stats, maxCount }: Props = $props();
  const id = $props.id();

  const pct = (count: number) => (maxCount > 0 ? (count / maxCount) * 100 : 0);
  const isDual = $derived(problem.holdColorSecondary !== undefined);
  const topPct = $derived(pct(stats.top + stats.flash));
</script>

<div class="container">
  <button {id} aria-label="Show details for problem #{problem.number}">
    <div
      class:dual={isDual}
      class="fill"
      style:--primary-color={problem.holdColorPrimary}
      style:--secondary-color={problem.holdColorSecondary ?? problem.holdColorPrimary}
    >
      {#if isDual}
        <div class="dual-top" style:--target-height="{topPct}%"></div>
      {:else}
        <div class="side" style:--color={problem.holdColorPrimary}>
          <SubBar percentage={topPct} fillWeight={1} />
        </div>
      {/if}
    </div>
  </button>
  <span class="label">#{problem.number}</span>

  <wa-popover for={id} placement="top">
    <strong>Problem #{problem.number}</strong>

    {#if problem.zone1Enabled}
      <div>
        Zone 1: {stats.zone1 + stats.zone2 + stats.top + stats.flash}
      </div>
    {/if}

    {#if problem.zone2Enabled}
      <div>Zone 2: {stats.zone2 + stats.top + stats.flash}</div>
    {/if}

    <div>Tops: {stats.top + stats.flash}</div>
    <div>Flashes: {stats.flash}</div>
  </wa-popover>
</div>

<style>
  .container {
    display: flex;
    flex-direction: column;
    width: 2rem;
    flex-shrink: 0;
  }

  button {
    overflow: hidden;
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 8rem;
    border: none;
    padding: 0;
    cursor: pointer;
    border-radius: var(--wa-border-radius-s);
    box-shadow: var(--wa-shadow-s);
  }

  .fill {
    overflow: hidden;
    position: relative;
    display: flex;
    height: 100%;
    background: rgb(from var(--primary-color) r g b / 2%);
    animation: grow var(--wa-transition-slow) var(--wa-transition-easing)
      forwards;
    transform: scaleY(0);
    transform-origin: bottom;
  }

  .fill.dual::before {
    content: "";
    position: absolute;
    inset: -50%;
    background: linear-gradient(
      to bottom,
      rgb(from var(--primary-color) r g b / 2%) 0 50%,
      rgb(from var(--secondary-color) r g b / 2%) 50% 100%
    );
    transform: rotate(-45deg);
  }

  .dual-top {
    position: absolute;
    inset-inline: 0;
    inset-block-end: 0;
    z-index: 1;
    height: var(--target-height);
    background: linear-gradient(
      135deg,
      var(--primary-color) 0 50%,
      var(--secondary-color) 50% 100%
    );
  }

  .side {
    position: relative;
    z-index: 1;
    display: flex;
    flex: 1;
    flex-direction: column-reverse;
  }

  .label {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
    margin-block-start: var(--wa-space-xs);
    text-align: center;
  }

  wa-popover::part(body) {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-3xs);
    font-size: var(--wa-font-size-s);
  }

  @keyframes grow {
    to {
      transform: scaleY(1);
    }
  }
</style>
