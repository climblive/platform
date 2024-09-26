<script lang="ts">
  import type { ScoreboardEntry } from "@climblive/lib/models";
  import "@shoelace-style/shoelace/dist/components/progress-bar/progress-bar.js";
  import { getContext } from "svelte";
  import { type Readable } from "svelte/store";
  import Floater from "./Floater.svelte";
  import ResultEntry from "./ResultEntry.svelte";

  export let compClassId: number;

  const scoreboard =
    getContext<Readable<Map<number, ScoreboardEntry[]>>>("scoreboard");
  $: results = $scoreboard.get(compClassId);
</script>

{#if results}
  <div
    style="height: calc({results.length} * 2.25rem + {results.length -
      1} * var(--sl-spacing-x-small))"
  >
    {#each results as scoreboardEntry (scoreboardEntry.contenderId)}
      <Floater order={scoreboardEntry.rankOrder}>
        <ResultEntry {scoreboardEntry} />
      </Floater>
    {/each}
  </div>
{:else}
  <sl-progress-bar indeterminate></sl-progress-bar>
{/if}

<style>
  div {
    position: relative;
    overflow: hidden;
    width: 100%;
  }

  sl-progress-bar {
    --height: 4px;
  }
</style>
