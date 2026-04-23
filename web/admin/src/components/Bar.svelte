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

  const topPct = $derived(pct(stats.top));
  const flashPct = $derived(pct(stats.flash));
</script>

<div class="container">
  <button {id} aria-label="Show details for problem #{problem.number}">
    <div class="fill" style:--color={problem.holdColorPrimary}>
      <SubBar percentage={flashPct} fillWeight={1} />
      <SubBar percentage={topPct} fillWeight={0.6} />
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
    display: flex;
    flex-direction: column-reverse;
    height: 100%;
    background: rgb(from var(--color) r g b / 2%);
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
</style>
