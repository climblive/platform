<script lang="ts">
  import type { ScoreboardEntry } from "@climblive/lib/models";
  import "@shoelace-style/shoelace/dist/components/progress-bar/progress-bar.js";
  import { getContext, onDestroy, onMount } from "svelte";
  import { type Readable } from "svelte/store";
  import Floater from "./Floater.svelte";
  import ResultEntry from "./ResultEntry.svelte";

  export let compClassId: number;

  let container: HTMLDivElement | undefined;
  let observer: ResizeObserver | undefined;
  let timerId: number;

  let height: number;
  let pageSize: number;
  let pageIndex = 0;

  $: {
    pageSize = Math.floor(height / 42);
    pageIndex = 0;
  }

  $: pageCount = Math.ceil(results.length / pageSize);

  const scoreboard =
    getContext<Readable<Map<number, ScoreboardEntry[]>>>("scoreboard");
  $: results = $scoreboard.get(compClassId) ?? [];

  onMount(() => {
    timerId = setInterval(() => {
      pageIndex = (pageIndex + 1) % pageCount;
    }, 10_000);

    observer = new ResizeObserver((entries) => {
      for (const entry of entries) {
        if (entry.contentBoxSize) {
          const contentBoxSize = entry.contentBoxSize[0];
          height = contentBoxSize.blockSize;
        }
      }
    });
  });

  onDestroy(() => {
    observer?.disconnect();
    observer = undefined;

    clearInterval(timerId);
    timerId = 0;
  });

  $: {
    if (container) {
      observer?.observe(container);
    }
  }
</script>

<div
  class="container"
  bind:this={container}
  style="--page-size: {pageSize}; --page-index: {pageIndex}"
>
  {#each results as scoreboardEntry (scoreboardEntry.contenderId)}
    <Floater order={scoreboardEntry.rankOrder}>
      <ResultEntry {scoreboardEntry} />
    </Floater>
  {/each}
</div>

<style>
  .container {
    overflow: hidden;
    position: relative;
    height: 100%;
    clip-path: rect(
      0px 0px
        calc(
          var(--page-size) * (2.25rem) + (var(--page-size) - 1) *
            var(--sl-spacing-x-small)
        )
        100%
    );
  }
</style>
