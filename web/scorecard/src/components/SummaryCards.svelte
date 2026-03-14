<script lang="ts">
  import { Timer } from "@climblive/lib/components";
  import type { Tick } from "@climblive/lib/models";
  import { type ContestState } from "@climblive/lib/types";
  import { ordinalSuperscript } from "@climblive/lib/utils";

  interface Props {
    score: number;
    placement: number | undefined;
    disqualified: boolean;
    contestState: ContestState;
    startTime: Date;
    endTime: Date;
    ticks: Tick[];
    totalProblems: number;
  }

  let {
    score,
    placement,
    disqualified,
    contestState,
    startTime,
    endTime,
    ticks,
    totalProblems,
  }: Props = $props();

  const formatScore = (value: number) => {
    return value.toLocaleString("en", { useGrouping: true }).replace(/,/g, " ");
  };

  const tops = $derived(ticks.filter((tick) => tick.top).length);
  const flashes = $derived(
    ticks.filter((tick) => tick.top && tick.attemptsTop === 1).length,
  );
</script>

<div class="cards">
  <div class="card">
    <span class="label">Points</span>
    <span class="big-value">{formatScore(score)}</span>
    <span class="label">Placement</span>
    <span class="value">
      {#if disqualified}
        Disqualified
      {:else if placement}
        {placement}<sup>{ordinalSuperscript(placement)}</sup>
      {:else}
        -
      {/if}
    </span>
  </div>
  <div class="card">
    {#if contestState === "NOT_STARTED"}
      <Timer endTime={startTime} label="Time until start" />
    {:else}
      <Timer {endTime} label="Time remaining" />
    {/if}
    <span class="label">Cleared</span>
    <span class="value">
      {tops}/{totalProblems} Problems
      {#if flashes > 0}
        <span class="flashes"
          >({flashes} {flashes === 1 ? "Flash" : "Flashes"})</span
        >
      {/if}
    </span>
  </div>
</div>

<style>
  .cards {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--wa-space-s);
  }

  .card {
    background-color: var(--wa-color-surface-raised);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-surface-border);
    border-radius: var(--wa-border-radius-m);
    padding: var(--wa-space-m);
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-3xs);
  }

  .label {
    font-size: var(--wa-font-size-2xs);
    font-weight: var(--wa-font-weight-bold);
    text-transform: uppercase;
    color: var(--wa-color-text-quiet);
    letter-spacing: 0.05em;
  }

  .big-value {
    font-size: var(--wa-font-size-3xl);
    font-weight: var(--wa-font-weight-bold);
    line-height: var(--wa-line-height-condensed);
    margin-bottom: var(--wa-space-xs);
  }

  .value {
    font-size: var(--wa-font-size-m);
    font-weight: var(--wa-font-weight-bold);
    line-height: var(--wa-line-height-condensed);
  }

  .flashes {
    color: var(--wa-color-brand-text-normal);
    font-weight: var(--wa-font-weight-normal);
  }

  .card :global(.timer) {
    margin-bottom: var(--wa-space-xs);
  }

  .card :global(.timer span[role="timer"]) {
    font-size: var(--wa-font-size-3xl);
    font-weight: var(--wa-font-weight-bold);
    line-height: var(--wa-line-height-condensed);
  }

  .card :global(.timer label) {
    font-size: var(--wa-font-size-2xs);
    font-weight: var(--wa-font-weight-bold);
    text-transform: uppercase;
    color: var(--wa-color-text-quiet);
    letter-spacing: 0.05em;
    order: -1;
  }
</style>
