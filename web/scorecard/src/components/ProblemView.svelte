<script lang="ts">
  import { Score } from "@climblive/lib/components";
  import type { Problem, Tick } from "@climblive/lib/models";
  import { calculateProblemScore } from "@climblive/lib/utils";
  import HoldColorIndicator from "./HoldColorIndicator.svelte";
  import TickBox from "./TickBox.svelte";

  interface Props {
    problem: Problem;
    tick?: Tick | undefined;
    disabled: boolean;
    maxProblemNumber: number;
  }

  let { problem, tick, disabled, maxProblemNumber }: Props = $props();

  let pointValue = $derived(calculateProblemScore(problem, tick));
</script>

<section
  data-ticked={!!tick}
  data-flashed={tick?.attemptsTop === 1}
  aria-label={`Problem ${problem.number}`}
  style="--number-length: {maxProblemNumber.toString().length + 2}ch"
>
  <span class="number">№ {problem.number}</span>
  <HoldColorIndicator
    primary={problem.holdColorPrimary}
    secondary={problem.holdColorSecondary}
  />
  <span class="points">
    <span class="top">
      {problem.pointsTop}p
    </span>
    {#if problem.flashBonus}
      <sl-icon name="lightning-charge"></sl-icon>
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
    background-color: var(--sl-color-stone-50);
    border-radius: var(--sl-border-radius-small);
    border: solid 1px var(--sl-color-stone-300);
    padding-left: var(--sl-spacing-small);
    padding-right: var(--sl-spacing-2x-small);

    display: grid;
    grid-template-columns: var(--number-length) 1rem 1fr 1fr 2.5rem;
    grid-template-rows: 1fr;
    gap: var(--sl-spacing-x-small);
    align-items: center;
    justify-items: end;
  }

  .number {
    font-size: var(--sl-font-size-small);
    text-wrap: nowrap;
    justify-self: start;
    font-weight: var(--sl-font-weight-bold);
  }

  .points {
    margin-right: auto;
    white-space: nowrap;

    & sl-icon {
      font-size: var(--sl-font-size-x-small);
      color: var(--sl-color-yellow-500);
    }
  }

  .score {
    text-align: right;
  }
</style>
