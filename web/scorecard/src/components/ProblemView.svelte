<script lang="ts">
  import { HoldColorIndicator, Score } from "@climblive/lib/components";
  import type { PointValue, Problem, Tick } from "@climblive/lib/models";
  import TickDialog from "./TickDialog.svelte";

  type ScorecardProblem = Problem & {
    pointValue?: PointValue;
  };

  interface Props {
    problem: ScorecardProblem;
    tick?: Tick | undefined;
    disabled: boolean;
    counted: boolean;
    enablePoints: boolean;
    enableAttempts: boolean;
  }

  const {
    problem,
    tick,
    disabled,
    counted,
    enablePoints,
    enableAttempts,
  }: Props = $props();

  const id = $props.id();

  const valueRange = $derived.by<{ min: number; max: number } | undefined>(
    () => {
      if (problem.pointValue === undefined) {
        return undefined;
      }

      const { zone1, zone2, top, flashBonus } = problem.pointValue;

      const values = [zone1, zone2, top + flashBonus];

      return { min: Math.min(...values), max: Math.max(...values) };
    },
  );
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
  <span class="number">#{problem.number}</span>
  <span class="points">
    {#if enablePoints && valueRange}
      <span class="top">
        {#if valueRange.min === valueRange.max}
          <Score value={valueRange.max + "p"} compactThreshold={9_999_999}
          ></Score>
        {:else}
          <span class="lte">≤</span>
          <Score value={valueRange.max + "p"} compactThreshold={9_999_999}
          ></Score>
        {/if}
      </span>
    {/if}

    {#if problem.description}
      <wa-popover for={id}>
        {problem.description}
        <wa-button
          appearance="outlined"
          variant="neutral"
          size="s"
          data-popover="close">Got it!</wa-button
        >
      </wa-popover>

      <wa-button size="xs" appearance="plain" pill {id}
        ><wa-icon name="info-circle"></wa-icon></wa-button
      >
    {/if}
  </span>
  <div class="score" class:uncounted={!counted}>
    {#if enablePoints && tick && problem.pointValue?.current !== undefined}
      <Score
        value={problem.pointValue.current + "p"}
        prefix={counted ? "+" : undefined}
        compactThreshold={9_999}
      />
    {:else if tick?.top}
      +1t
    {:else if tick?.zone2 && problem.zone2Enabled}
      +1z₂
    {:else if tick?.zone1 && problem.zone1Enabled}
      +1z₁
    {/if}
  </div>

  <TickDialog
    {problem}
    {tick}
    {disabled}
    pointValue={problem.pointValue}
    {enableAttempts}
    {enablePoints}
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
    display: flex;
    align-items: center;
    gap: var(--wa-space-2xs);

    & .top {
      white-space: nowrap;
    }

    & .lte {
      color: var(--wa-color-text-quiet);
    }
  }

  .score {
    text-align: right;
    font-weight: var(--wa-font-weight-semibold);

    &.uncounted {
      opacity: 0.3;
      text-decoration: line-through;
    }
  }

  wa-popover::part(body) {
    width: 300px;
  }

  wa-popover wa-button {
    margin-block-start: var(--wa-space-m);
  }
</style>
