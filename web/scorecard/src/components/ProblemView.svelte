<script lang="ts">
  import { HoldColorIndicator, Score } from "@climblive/lib/components";
  import type { PointValue, Problem, Tick } from "@climblive/lib/models";
  import TickBox from "./TickBox.svelte";

  type ScorecardProblem = Problem & {
    pointValue?: PointValue;
  };

  interface Props {
    problem: ScorecardProblem;
    tick?: Tick | undefined;
    disabled: boolean;
    counted: boolean;
  }

  const { problem, tick, disabled, counted }: Props = $props();

  const valueRange = $derived.by<{ min: number; max: number }>(() => {
    if (problem.pointValue === undefined) {
      return { min: 0, max: 0 };
    }

    const { zone1, zone2, top, flash } = problem.pointValue;

    const values = [zone1, zone2, top, flash];

    return { min: Math.min(...values), max: Math.max(...values) };
  });
</script>

<section
  data-ticked={!!tick}
  data-flashed={tick?.attemptsTop === 1}
  aria-label={`Problem ${problem.number}`}
>
  <HoldColorIndicator
    primary={problem.holdColorPrimary}
    secondary={problem.holdColorSecondary}
    --height="1.25rem"
    --width="1.25rem"
  />
  <span class="number">№ {problem.number}</span>
  <span class="points">
    <span class="top">
      {#if valueRange.min === valueRange.max}
        {valueRange.max}p
      {:else}
        ≤ {valueRange.max}p
      {/if}
    </span>
  </span>
  <div class="score" class:uncounted={!counted}>
    {#if tick && problem.pointValue !== undefined}
      <Score
        value={problem.pointValue.current}
        prefix={counted ? "+" : undefined}
      />
    {/if}
  </div>

  <TickBox {problem} {tick} {disabled} pointValue={problem.pointValue} />
</section>

<style>
  section {
    height: 3rem;
    background-color: var(--wa-color-surface-raised);
    border-radius: var(--wa-border-radius-m);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-surface-border);
    padding-inline-start: var(--wa-space-s);
    padding-inline-end: var(--wa-space-xs);

    display: grid;
    grid-template-columns: max-content max-content 1fr 1fr 2.5rem;
    grid-column: 1 / -1;
    gap: var(--wa-space-xs);
    align-items: center;
    justify-items: end;
  }

  @supports (grid-template-columns: subgrid) {
    section {
      grid-template-columns: subgrid;
    }
  }

  .number {
    font-size: var(--wa-font-size-s);
    text-wrap: nowrap;
    justify-self: start;
    font-weight: var(--wa-font-weight-bold);
  }

  .points {
    margin-right: auto;
    white-space: nowrap;
  }

  .score {
    text-align: right;

    &.uncounted {
      opacity: 0.3;
      text-decoration: line-through;
    }
  }
</style>
