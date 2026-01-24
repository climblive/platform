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
  <div class="row">
    <div class="stat">
      <span class="label">Tops</span>
      <span class="value"><strong>{tops}</strong>/{totalProblems}</span>
    </div>
    {#if hasZones}
      <div class="stat">
        <span class="label">Zones</span>
        <span class="value"><strong>{zones}</strong>/{totalProblems}</span>
      </div>
    {/if}
    <div class="stat">
      <span class="label">Flashes</span>
      <span class="value"><strong>{flashes}</strong>/{totalProblems}</span>
    </div>
  </div>
  <div class="row">
    <div class="stat">
      <span class="label">Score</span>
      <span class="value"><strong>{score}</strong></span>
    </div>
    <div class="stat">
      <span class="label">Placement</span>
      <span class="value">
        {#if placement}
          <strong>{placement}</strong><sup>{ordinalSuperscript(placement)}</sup>
        {:else}
          <strong>-</strong>
        {/if}
      </span>
    </div>
  </div>
</div>

<style>
  .summary {
    background-color: var(--wa-color-surface-raised);
    border: var(--wa-border-width-m) var(--wa-border-style)
      var(--wa-color-surface-border);
    border-radius: var(--wa-border-radius-l);
    padding: var(--wa-space-m);
    margin-bottom: var(--wa-space-m);
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }

  .row {
    display: flex;
    flex-wrap: wrap;
    gap: var(--wa-space-m);
  }

  .stat {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    min-width: 4rem;
  }

  .label {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-weak);
    margin-bottom: var(--wa-space-2xs);
  }

  .value {
    font-size: 1em;
    line-height: 1;

    & strong {
      font-size: 1.5em;
      font-weight: var(--wa-font-weight-bold);
    }
  }
</style>
