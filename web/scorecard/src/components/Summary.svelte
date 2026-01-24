<script lang="ts">
  import type { Problem, Tick } from "@climblive/lib/models";
  import { ordinalSuperscript } from "@climblive/lib/utils";

  interface Props {
    ticks: Tick[];
    problems: Problem[];
    score: number;
    placement: number | undefined;
  }

  let { ticks, problems, score, placement }: Props = $props();

  let tops = $derived(ticks.filter((tick) => tick.top).length);

  let zones = $derived(
    ticks.filter((tick) => !tick.top && (tick.zone1 || tick.zone2)).length,
  );

  let flashes = $derived(
    ticks.filter((tick) => tick.top && tick.attemptsTop === 1).length,
  );

  let hasZones = $derived(
    problems.some((problem) => problem.zone1Enabled || problem.zone2Enabled),
  );

  let totalProblems = $derived(problems.length);
</script>

<div class="summary">
  <h2>Your Results</h2>
  <div class="stats">
    <div class="stat">
      <span class="label">Tops</span>
      <span class="value">{tops}/{totalProblems}</span>
    </div>
    {#if hasZones}
      <div class="stat">
        <span class="label">Zones</span>
        <span class="value">{zones}/{totalProblems}</span>
      </div>
    {/if}
    <div class="stat">
      <span class="label">Flashes</span>
      <span class="value">{flashes}/{totalProblems}</span>
    </div>
    <div class="stat">
      <span class="label">Score</span>
      <span class="value">{score}</span>
    </div>
    <div class="stat">
      <span class="label">Placement</span>
      <span class="value">
        {#if placement}
          {placement}<sup>{ordinalSuperscript(placement)}</sup>
        {:else}
          -
        {/if}
      </span>
    </div>
  </div>
</div>

<style>
  .summary {
    background-color: var(--wa-color-surface-default);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-border-default);
    border-radius: var(--wa-border-radius-m);
    padding: var(--wa-space-m);
    margin-bottom: var(--wa-space-m);
  }

  h2 {
    margin: 0 0 var(--wa-space-s) 0;
    font-size: var(--wa-font-size-m);
    font-weight: var(--wa-font-weight-bold);
  }

  .stats {
    display: flex;
    flex-wrap: wrap;
    gap: var(--wa-space-m);
  }

  .stat {
    display: flex;
    flex-direction: column;
    align-items: center;
    min-width: 4rem;
  }

  .label {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-weak);
    margin-bottom: var(--wa-space-2xs);
  }

  .value {
    font-size: var(--wa-font-size-l);
    font-weight: var(--wa-font-weight-bold);
  }
</style>
