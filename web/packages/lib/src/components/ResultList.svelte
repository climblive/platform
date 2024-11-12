<script lang="ts">
  import type { ScoreboardEntry } from "@climblive/lib/models";
  import "@shoelace-style/shoelace/dist/components/progress-bar/progress-bar.js";
  import { getContext, onDestroy, onMount } from "svelte";
  import { type Readable } from "svelte/store";
  import Floater from "./Floater.svelte";
  import ResultEntry from "./ResultEntry.svelte";

  export let compClassId: number;
  export let overflow: "pagination" | "scroll" = "scroll";

  const ITEM_HEIGHT = 36;
  const GAP = 8;

  let container: HTMLDivElement | undefined;
  let observer: ResizeObserver | undefined;
  let pageFlipIntervalTimerId: number;

  let containerHeight: number;
  let pageSize: number;
  let pageIndex = 0;

  $: {
    switch (overflow) {
      case "pagination":
        pageSize = Math.floor((containerHeight + GAP) / (ITEM_HEIGHT + GAP));
        break;
      case "scroll":
        pageSize = results.length;
        break;
    }
  }

  $: {
    if (overflow === "scroll") {
      pageIndex = 0;
    }
  }

  $: pageCount = Math.ceil(results.length / pageSize);

  const scoreboard =
    getContext<Readable<Map<number, ScoreboardEntry[]>>>("scoreboard");
  $: results = $scoreboard.get(compClassId) ?? [];

  onMount(() => {
    pageFlipIntervalTimerId = setInterval(() => {
      if (overflow === "pagination") {
        pageIndex = (pageIndex + 1) % pageCount;
      }
    }, 10_000);

    observer = new ResizeObserver((entries) => {
      for (const entry of entries) {
        if (entry.contentBoxSize) {
          const contentBoxSize = entry.contentBoxSize[0];
          containerHeight = contentBoxSize.blockSize;
        }
      }
    });
  });

  onDestroy(() => {
    observer?.disconnect();
    observer = undefined;

    clearInterval(pageFlipIntervalTimerId);
    pageFlipIntervalTimerId = 0;
  });

  $: {
    if (container && overflow === "pagination") {
      observer?.observe(container);
    }
  }
</script>

<div
  class="container"
  bind:this={container}
  style="--page-size: {pageSize}; --page-index: {pageIndex}"
  {overflow}
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
  }

  .container[overflow="pagination"] {
    height: 100%;
    clip-path: rect(
      0px 100%
        calc(
          var(--page-size) * 2.25rem + (var(--page-size) - 1) *
            var(--sl-spacing-x-small)
        )
        0px
    );
  }

  .container[overflow="scroll"] {
    height: calc(var(--page-size) * (2.25rem + var(--sl-spacing-x-small)));
  }
</style>
