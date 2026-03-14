<script lang="ts">
  import { Timer } from "@climblive/lib/components";
  import { type ContestState } from "@climblive/lib/types";
  import { ordinalSuperscript } from "@climblive/lib/utils";

  interface Props {
    score: number;
    placement: number | undefined;
    disqualified: boolean;
    contestState: ContestState;
    startTime: Date;
    endTime: Date;
  }

  let {
    score,
    placement,
    disqualified,
    contestState,
    startTime,
    endTime,
  }: Props = $props();

  const formatScore = (value: number) => {
    return value.toLocaleString("en", { useGrouping: true }).replace(/,/g, " ");
  };

</script>

<div class="card">
  <div class="card-section">
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
  <div class="card-section">
    {#if contestState === "NOT_STARTED"}
      <Timer endTime={startTime} label="Time until start" />
    {:else}
      <Timer {endTime} label="Time remaining" />
    {/if}
  </div>
</div>

<style>
  .card {
    background-color: var(--wa-color-surface-raised);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-surface-border);
    border-radius: var(--wa-border-radius-m);
    padding: var(--wa-space-m);
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--wa-space-m);
  }

  .card-section {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-3xs);
  }

  .label {
    font-size: var(--wa-font-size-2xs);
    font-weight: var(--wa-font-weight-bold);
    color: var(--wa-color-text-quiet);
  }

  .big-value {
    font-size: var(--wa-font-size-xl);
    font-weight: var(--wa-font-weight-bold);
    line-height: var(--wa-line-height-condensed);
    margin-bottom: var(--wa-space-xs);
  }

  .value {
    font-size: var(--wa-font-size-m);
    font-weight: var(--wa-font-weight-bold);
    line-height: var(--wa-line-height-condensed);
  }

  .card :global(.timer) {
    margin-bottom: var(--wa-space-xs);
  }

  .card :global(.timer span[role="timer"]) {
    font-size: var(--wa-font-size-xl);
    font-weight: var(--wa-font-weight-bold);
    line-height: var(--wa-line-height-condensed);
  }

  .card :global(.timer label) {
    font-size: var(--wa-font-size-2xs);
    font-weight: var(--wa-font-weight-bold);
    color: var(--wa-color-text-quiet);
    order: -1;
  }
</style>
