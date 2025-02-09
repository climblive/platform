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

  let padding = $derived.by(() => {
    const length =
      maxProblemNumber.toString().length - problem.number.toString().length;

    if (length <= 0) {
      return "";
    }

    let padding = "";

    for (let k = 0; k < length; k++) {
      padding += "0";
    }

    return padding;
  });
</script>

<section
  data-ticked={!!tick}
  data-flashed={tick?.attemptsTop === 1}
  aria-label={`Problem ${problem.number}`}
>
  <span class="number"
    >№ <span class="padding">{padding}</span>{problem.number}</span
  >
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
    background-color: var(--sl-color-primary-50);
    border-radius: var(--sl-border-radius-small);
    border: solid 1px var(--sl-color-primary-300);
    padding-left: var(--sl-spacing-small);
    padding-right: var(--sl-spacing-2x-small);
    color: var(--sl-color-primary-950);

    display: grid;
    grid-template-columns: max-content 1rem 1fr 1fr 2.5rem;
    grid-template-rows: 1fr;
    gap: var(--sl-spacing-x-small);
    align-items: center;
    justify-items: end;
  }

  .number {
    font-size: var(--sl-font-size-x-small);
    text-wrap: nowrap;
  }

  .padding {
    visibility: hidden;
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
