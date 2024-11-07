<script lang="ts">
  import { Timer } from "@climblive/lib/components";
  import type { ScoreboardEntry } from "@climblive/lib/models";
  import { useContestState } from "@climblive/lib/utils";
  import { getContext, onDestroy } from "svelte";
  import { type Readable } from "svelte/store";

  export let name: string;
  export let compClassId: number;
  export let startTime: Date;
  export let endTime: Date;
  export let gracePeriodEndTime: Date;

  const { state, stop, update } = useContestState();

  $: {
    if (startTime && endTime && gracePeriodEndTime) {
      update(startTime, endTime, gracePeriodEndTime);
    }
  }

  onDestroy(() => {
    stop();
  });

  const scoreboard =
    getContext<Readable<Map<number, ScoreboardEntry[]>>>("scoreboard");
  $: results = $scoreboard.get(compClassId) ?? [];

  $: allContenders = [...$scoreboard.values()].reduce((count, results) => {
    return count + results.length;
  }, 0);
</script>

<header>
  <h2>{name} <span class="size">({results.length}/{allContenders})</span></h2>
  <div class="timer">
    {#if $state === "NOT_STARTED"}
      <Timer endTime={startTime} />
      <span class="footer">Time until start</span>
    {:else}
      <Timer {endTime} />
      <span class="footer">Time remaining</span>
    {/if}
  </div>
</header>

<style>
  header {
    margin-bottom: var(--sl-spacing-large);

    display: flex;
    flex-direction: column;
    align-items: center;

    & h2 {
      line-height: var(--sl-line-height-denser);
    }
  }

  .size {
    font-size: var(--sl-font-size-small);
    color: var(--sl-color-primary-700);
  }

  .timer {
    text-align: center;
    font-weight: var(--sl-font-weight-bold);
    font-size: var(--sl-font-size-large);

    & > * {
      display: block;
    }

    & .footer {
      font-weight: var(--sl-font-weight-normal);
      font-size: var(--sl-font-size-small);
    }
  }
</style>
