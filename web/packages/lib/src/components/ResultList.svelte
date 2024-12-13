<script lang="ts">
  import { run } from 'svelte/legacy';

  import type { ScoreboardEntry } from "@climblive/lib/models";
  import "@shoelace-style/shoelace/dist/components/progress-bar/progress-bar.js";
  import "@shoelace-style/shoelace/dist/components/skeleton/skeleton.js";
  import { onDestroy, onMount } from "svelte";
  import { type Readable } from "svelte/store";
  import Floater from "./Floater.svelte";
  import ResultEntry from "./ResultEntry.svelte";

  interface Props {
    compClassId: number;
    overflow?: "pagination" | "scroll";
    scoreboard: Readable<Map<number, ScoreboardEntry[]>>;
    loading: boolean;
  }

  let {
    compClassId,
    overflow = "scroll",
    scoreboard,
    loading
  }: Props = $props();

  const ITEM_HEIGHT = 36;
  const GAP = 8;
  const SCROLLABLE_SKELETON_ENTRIES = 10;

  let container: HTMLDivElement | undefined = $state();
  let observer: ResizeObserver | undefined;
  let pageFlipIntervalTimerId: number;

  let containerHeight: number = $state(0);
  let pageSize: number = $state(0);
  let pageIndex = $state(0);





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

  let results = $derived($scoreboard.get(compClassId) ?? []);
  run(() => {
    switch (overflow) {
      case "pagination":
        pageSize = Math.floor((containerHeight + GAP) / (ITEM_HEIGHT + GAP));
        break;
      case "scroll":
        pageSize = results.length;

        if (loading) {
          pageSize = SCROLLABLE_SKELETON_ENTRIES;
        }

        break;
    }
  });
  run(() => {
    if (overflow === "scroll") {
      pageIndex = 0;
    }
  });
  let pageCount = $derived(Math.ceil(results.length / pageSize));
  run(() => {
    if (container && overflow === "pagination") {
      observer?.observe(container);
    }
  });
</script>

<div
  class="container"
  bind:this={container}
  style="--page-size: {pageSize}; --page-index: {pageIndex}"
  {overflow}
>
  {#if loading}
    {#each [...Array(overflow === "pagination" ? pageSize : SCROLLABLE_SKELETON_ENTRIES).keys()] as i (i)}
      <sl-skeleton effect="sheen"></sl-skeleton>
    {/each}
  {:else}
    {#each results as scoreboardEntry (scoreboardEntry.contenderId)}
      {#if scoreboardEntry.score}
        <Floater order={scoreboardEntry.score.rankOrder}>
          <ResultEntry {scoreboardEntry} />
        </Floater>
      {/if}
    {/each}
  {/if}
</div>

<style>
  sl-skeleton {
    --color: var(--sl-color-primary-400);
    --sheen-color: var(--sl-color-primary-300);
    --border-radius: var(--sl-border-radius-medium);

    margin-bottom: var(--sl-spacing-x-small);
    height: 2.25rem;
  }

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
