<script lang="ts">
  import { Timer } from "@climblive/lib/components";
  import { type ContestState } from "@climblive/lib/types";
  import { ordinalSuperscript } from "@climblive/lib/utils";
  import type { Snippet } from "svelte";

  interface Props {
    score: number;
    placement: number | undefined;
    disqualified: boolean;
    contestState: ContestState;
    startTime: Date;
    endTime: Date;
  }

  const {
    score,
    placement,
    disqualified,
    contestState,
    startTime,
    endTime,
  }: Props = $props();
</script>

{#snippet entry(label: string, value: Snippet, disabled: boolean = false)}
  <div class="stat" class:disabled>
    <span class="label">{label}</span>
    <span class="value">{@render value()}</span>
  </div>
{/snippet}

{#snippet pointsValue()}
  <strong>{score}p</strong>
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

{#snippet timerValue()}
  {#if contestState === "NOT_STARTED"}
    <Timer endTime={startTime} />
  {:else}
    <Timer {endTime} />
  {/if}
{/snippet}

<div class="summary">
  {@render entry("Score", pointsValue)}
  {@render entry("Placement", placementValue)}

  {#if contestState === "NOT_STARTED"}
    {@render entry("Time until start", timerValue)}
  {:else}
    {@render entry("Time remaining", timerValue)}
  {/if}
</div>

<style>
  .summary {
    background-color: var(--wa-color-surface-raised);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-surface-border);
    border-radius: var(--wa-border-radius-m);
    padding: var(--wa-space-m);
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
  }

  .stat.disabled {
    opacity: 0.5;
  }

  .summary :global(.timer) {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-2xs);
  }

  .summary :global(.timer label) {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
  }

  .summary :global(.timer span[role="timer"]) {
    font-size: 1.5em;
    font-weight: var(--wa-font-weight-bold);
    line-height: 1;
  }
</style>
