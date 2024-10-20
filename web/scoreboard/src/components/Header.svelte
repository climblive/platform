<script lang="ts">
  import { Timer } from "@climblive/lib/components";
  import type { ScoreboardEntry } from "@climblive/lib/models";
  import { getContext } from "svelte";
  import { type Readable } from "svelte/store";

  export let name: string;
  export let compClassId: number;
  export let timeEnd: Date;

  const scoreboard =
    getContext<Readable<Map<number, ScoreboardEntry[]>>>("scoreboard");
  $: results = $scoreboard.get(compClassId) ?? [];

  $: allContenders = [...$scoreboard.values()].reduce((count, results) => {
    return count + results.length;
  }, 0);
</script>

<header>
  <h1>{name} <span class="size">({results.length}/{allContenders})</span></h1>
  <Timer endTime={timeEnd} />
</header>

<style>
  header {
    text-align: center;
    margin-bottom: var(--sl-spacing-large);

    & h1 {
      line-height: var(--sl-line-height-denser);
    }
  }

  .size {
    font-size: var(--sl-font-size-small);
    color: var(--sl-color-primary-700);
  }
</style>
