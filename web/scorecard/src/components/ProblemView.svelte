<script lang="ts">
  import { Score } from "@climblive/lib/components";
  import type { Problem, Tick } from "@climblive/lib/models";
  import { calculateProblemScore } from "@climblive/lib/utils";
  import HoldColorIndicator from "./HoldColorIndicator.svelte";
  import TickBox from "./TickBox.svelte";

  export let problem: Problem;
  export let tick: Tick | undefined = undefined;
  export let disabled: boolean;

  $: pointValue = calculateProblemScore(problem, tick);
</script>

<section data-ticked={!!tick} data-flashed={tick?.attemptsTop === 1}>
  <span class="number">{problem.number}.</span>
  <HoldColorIndicator
    primary={problem.holdColorPrimary}
    secondary={problem.holdColorSecondary}
  />
  <span class="points">
    <span class="top">
      {problem.pointsTop}p
    </span>
    {#if problem.flashBonus}
      <span class="flash">
        {problem.flashBonus}p
        <sl-icon name="lightning-charge"></sl-icon>
      </span>
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
    background-color: var(--sl-color-primary-100);
    border-radius: var(--sl-border-radius-small);
    border: solid 1px
      color-mix(in srgb, var(--sl-color-primary-300), transparent 50%);
    padding-left: var(--sl-spacing-small);
    padding-right: var(--sl-spacing-2x-small);
    color: var(--sl-color-primary-900);
    display: grid;
    grid-template-columns: 1rem 1.25rem 1fr 1fr 2.5rem;
    grid-template-rows: 1fr;
    gap: var(--sl-spacing-x-small);
    align-items: center;
    justify-items: end;
  }

  .number {
    font-size: var(--sl-font-size-x-small);
  }

  .points {
    margin-right: auto;
    white-space: nowrap;

    & .top {
      font-weight: var(--sl-font-weight-semibold);
    }

    & .flash {
      font-size: var(--sl-font-size-small);

      & sl-icon {
        font-size: var(--sl-font-size-x-small);
        color: var(--sl-color-yellow-500);
      }
    }
  }

  .score {
    text-align: right;
  }
</style>
