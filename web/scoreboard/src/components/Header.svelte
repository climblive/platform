<script lang="ts">
  import { ContestStateProvider, Timer } from "@climblive/lib/components";
  import type { ScoreboardEntry } from "@climblive/lib/models";
  import { getContext } from "svelte";
  import { type Readable } from "svelte/store";

  export let name: string;
  export let compClassId: number;
  export let startTime: Date;
  export let endTime: Date;

  const scoreboard =
    getContext<Readable<Map<number, ScoreboardEntry[]>>>("scoreboard");
  $: results = $scoreboard.get(compClassId) ?? [];

  $: allContenders = [...$scoreboard.values()].reduce((count, results) => {
    return count + results.length;
  }, 0);
</script>

<ContestStateProvider {startTime} {endTime} let:state>
  <header>
    <h2>{name} <span class="size">({results.length}/{allContenders})</span></h2>
    {#if state === "NOT_STARTED"}
      <Timer endTime={startTime} label="Time until start" />
    {:else}
      <Timer {endTime} label="Time remaining" />
    {/if}
  </header>
</ContestStateProvider>

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

  @media screen and (max-width: 512px) {
    header > h2 {
      display: none;
    }
  }

  .size {
    font-size: var(--sl-font-size-small);
    color: var(--sl-color-primary-700);
  }
</style>
