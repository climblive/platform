<script lang="ts">
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
    loading,
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
        if (pageCount === 0) {
          pageIndex = 0;
        } else {
          pageIndex = (pageIndex + 1) % pageCount;
        }
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

  $effect(() => {
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

  $effect(() => {
    if (overflow === "scroll") {
      pageIndex = 0;
    }
  });

  let pageCount = $derived(
    pageSize === 0 ? 0 : Math.ceil(results.length / pageSize),
  );

  $effect(() => {
    if (container && overflow === "pagination") {
      observer?.observe(container);
    }
  });
</script>

<div
  class="container"
  bind:this={container}
  style="--page-size: {pageSize}; --page-index: {pageIndex}"
  data-overflow={overflow}
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
{#if overflow === "pagination"}
  <div class="pagination">
    <div class="inset">
      {#each [...Array(pageCount).keys()] as i (i)}
        <div class="pageIndicator" data-current={i === pageIndex}></div>
      {/each}
    </div>
  </div>
{/if}

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

  .container[data-overflow="pagination"] {
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

  .container[data-overflow="scroll"] {
    height: calc(var(--page-size) * (2.25rem + var(--sl-spacing-x-small)));
  }

  .pagination {
    width: 100%;
    height: 0.5rem;
    overflow: hidden;
    position: relative;
    flex-shrink: 0;
    margin-top: var(--sl-spacing-small);

    & .inset {
      position: absolute;
      inset: 0;

      display: flex;
      justify-content: center;
      gap: 0.5rem;
    }
  }

  .pageIndicator {
    width: 0.5rem;
    height: 0.5rem;
    border-radius: 50%;
    background-color: var(--sl-color-primary-600);
    opacity: 0.5;

    &[data-current="true"] {
      opacity: 1;
    }
  }
</style>
