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
    disqualified: boolean;
  }

  const { problem, tick, disabled, disqualified }: Props = $props();
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
    {#if problem.pointValue}
      <span class="top">≤ {problem.pointValue.maximum}p</span>
    {/if}
  </span>
  <div class="score">
    {#if tick && problem.pointValue}
      <Score value={disqualified ? 0 : problem.pointValue.current} prefix="+" />
    {/if}
  </div>

  <TickBox
    {problem}
    {tick}
    {disabled}
    currentTickValue={problem.pointValue?.current}
  />
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
  }
</style>
