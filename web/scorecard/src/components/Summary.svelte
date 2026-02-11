<script lang="ts">
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import type { Problem, Tick } from "@climblive/lib/models";
  import { ordinalSuperscript } from "@climblive/lib/utils";
  import type { Snippet } from "svelte";

  interface Props {
    ticks: Tick[];
    problems: Problem[];
    score: number;
    placement: number | undefined;
    finalist: boolean;
    disqualified: boolean;
  }

  const { ticks, problems, score, placement, finalist, disqualified }: Props =
    $props();

  const tops = $derived(ticks.filter((tick) => tick.top).length);

  const zones = $derived(
    ticks.filter((tick) => {
      const problem = problems.find(({ id }) => id === tick.problemId);

      if (!(problem?.zone1Enabled || problem?.zone2Enabled)) {
        return false;
      }

      return tick.zone1 || tick.zone2;
    }).length,
  );

  const flashes = $derived(
    ticks.filter((tick) => tick.top && tick.attemptsTop === 1).length,
  );

  const problemsWithZones = $derived(
    problems.filter((problem) => problem.zone1Enabled || problem.zone2Enabled),
  );

  const totalProblems = $derived(problems.length);
</script>

{#snippet entry(label: string, value: Snippet, disabled: boolean = false)}
  <div class="stat" class:disabled>
    <span class="label">{label}</span>
    <span class="value">{@render value()}</span>
  </div>
{/snippet}

{#snippet topsValue()}
  <strong>{tops}</strong>/{totalProblems}
{/snippet}

{#snippet zonesValue()}
  <strong>{zones}</strong>/{problemsWithZones.length}
{/snippet}

{#snippet flashesValue()}
  <strong>{flashes}</strong>/{totalProblems}
{/snippet}

{#snippet scoreValue()}
  <strong>{score}</strong> pts
{/snippet}

{#snippet placementValue()}
  {#if disqualified}
    <strong>Disqualified</strong>
  {:else if placement}
    <strong>{placement}<sup>{ordinalSuperscript(placement)}</sup></strong>
  {:else}
    <strong>-</strong>
  {/if}
{/snippet}

{#snippet finalistValue()}
  <wa-icon name={finalist ? "medal" : "minus"}></wa-icon>
{/snippet}

<div class="summary">
  {@render entry("Tops", topsValue)}
  {@render entry("Zones", zonesValue, problemsWithZones.length === 0)}
  {@render entry("Flashes", flashesValue)}
  {@render entry("Score", scoreValue)}
  {@render entry("Placement", placementValue)}
  {@render entry("Finalist", finalistValue)}
</div>

<style>
  .summary {
    background-color: var(--wa-color-surface-raised);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-surface-border);
    border-radius: var(--wa-border-radius-m);
    padding: var(--wa-space-m);
    margin-bottom: var(--wa-space-m);
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: var(--wa-space-m);
  }

  .stat {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
  }

  .label {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
    margin-bottom: var(--wa-space-2xs);
  }

  .value {
    font-size: 1em;
    line-height: 1;

    & strong {
      font-size: 1.5em;
      font-weight: var(--wa-font-weight-bold);
    }

    & wa-icon {
      font-size: 1.5em;
    }
  }

  .stat.disabled {
    opacity: 0.5;
  }
</style>
