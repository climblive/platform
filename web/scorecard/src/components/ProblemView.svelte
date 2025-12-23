<script lang="ts">
  import { HoldColorIndicator, Score } from "@climblive/lib/components";
  import type { Problem, Tick } from "@climblive/lib/models";
  import { calculateProblemScore } from "@climblive/lib/utils";
  import TickBox from "./TickBox.svelte";

  interface Props {
    problem: Problem;
    tick?: Tick | undefined;
    disabled: boolean;
    highestProblemNumber: number;
  }

  const { problem, tick, disabled, highestProblemNumber }: Props = $props();

  const pointValue = $derived(calculateProblemScore(problem, tick));
</script>

<section
  data-ticked={!!tick}
  data-flashed={tick?.attemptsTop === 1}
  aria-label={`Problem ${problem.number}`}
  style="--number-length: {highestProblemNumber.toString().length + 2}ch"
>
  <span class="number">№ {problem.number}</span>
  <HoldColorIndicator
    primary={problem.holdColorPrimary}
    secondary={problem.holdColorSecondary}
    --height="1.25rem"
    --width="1.25rem"
  />
  <span class="points">
    <span class="top">
      {problem.pointsTop}p
    </span>
    {#if problem.flashBonus}
      <wa-icon name="bolt"></wa-icon>
    {/if}
  </span>
  <div class="score">
    <Score value={pointValue} hideZero prefix="+" />
  </div>

  <TickBox {problem} {tick} {disabled} />
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
    grid-template-columns: var(--number-length) 1rem 1fr 1fr 2.5rem;
    grid-template-rows: 1fr;
    gap: var(--wa-space-xs);
    align-items: center;
    justify-items: end;
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

    & wa-icon {
      font-size: var(--wa-font-size-xs);
    }
  }

  .score {
    text-align: right;
  }
</style>
