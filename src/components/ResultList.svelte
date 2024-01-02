<script lang="ts">
  import { getContext } from "svelte";
  import { type Readable } from "svelte/store";
  import Floater from "./Floater.svelte";
  import ResultEntry from "./ResultEntry.svelte";
  import type { RankedContender } from "../types";
  import "@shoelace-style/shoelace/dist/components/progress-bar/progress-bar.js";

  export let compClassId: number;

  const scoreboard =
    getContext<Readable<Map<number, RankedContender[]>>>("scoreboard");
  $: results = $scoreboard.get(compClassId);
</script>

{#if results}
  <div
    style="height: calc({results.length} * 2.25rem + {results.length -
      1} * var(--sl-spacing-x-small))"
  >
    {#each results as contender (contender.contenderId)}
      <Floater order={contender.order}>
        <ResultEntry
          placement={contender.placement}
          finalist={contender.finalist}
          {contender}
        />
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
