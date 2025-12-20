<script lang="ts">
  import "@awesome.me/webawesome/dist/components/progress-bar/progress-bar.js";
  import "@awesome.me/webawesome/dist/components/skeleton/skeleton.js";
  import { onDestroy, onMount } from "svelte";
  import { type Readable } from "svelte/store";
  import type { ScoreboardEntry } from "../models";
  import Floater from "./Floater.svelte";
  import ResultEntry from "./ResultEntry.svelte";

  interface Props {
    compClassId: number;
    overflow?: "pagination" | "scroll";
    scoreboard: Readable<Map<number, ScoreboardEntry[]>>;
    loading: boolean;
    highlightedContenderId?: number;
  }

  let {
    compClassId,
    overflow = "scroll",
    scoreboard,
    loading,
    highlightedContenderId,
  }: Props = $props();

  const ITEM_HEIGHT = 36;
  const GAP = 8;
  const SCROLLABLE_SKELETON_ENTRIES = 10;

  let container: HTMLDivElement | undefined = $state();
  let observer: ResizeObserver | undefined;
  let pageFlipIntervalTimerId: number;
  let visibilityObserver: IntersectionObserver | undefined;
  let isVisible = $state(true);

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
    }, 7_000);

    observer = new ResizeObserver((entries) => {
      for (const entry of entries) {
        if (entry.contentBoxSize) {
          const contentBoxSize = entry.contentBoxSize[0];
          containerHeight = contentBoxSize.blockSize;
        }
      }
    });

    visibilityObserver = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting) {
            isVisible = true;
          }
        }
      },
      { threshold: 0.1 }
    );
  });

  onDestroy(() => {
    observer?.disconnect();
    observer = undefined;

    visibilityObserver?.disconnect();
    visibilityObserver = undefined;

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

  $effect(() => {
    if (container && visibilityObserver) {
      visibilityObserver.observe(container);
    }
  });

  $effect(() => {
    if (highlightedContenderId && results.length > 0 && container && isVisible) {
      setTimeout(() => {
        if (!container) return;
        const highlightedEntry = container.querySelector(
          `section[data-highlighted="true"]`,
        );
        if (highlightedEntry) {
          highlightedEntry.scrollIntoView({
            behavior: "smooth",
            block: "center",
          });
        }
      });
      isVisible = false;
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
      <wa-skeleton effect="sheen"></wa-skeleton>
    {/each}
  {:else}
    {#each results as scoreboardEntry (scoreboardEntry.contenderId)}
      {#if scoreboardEntry.score}
        <Floater order={scoreboardEntry.score.rankOrder}>
          <ResultEntry
            {scoreboardEntry}
            highlighted={highlightedContenderId === scoreboardEntry.contenderId}
          />
        </Floater>
      {/if}
    {/each}
  {/if}
</div>

<style>
  wa-skeleton {
    margin-bottom: var(--wa-space-xs);
    height: 2.25rem;
  }

  wa-skeleton::part(indicator) {
    border-radius: var(--wa-border-radius-m);
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
            var(--wa-space-xs)
        )
        0px
    );
  }

  .container[data-overflow="scroll"] {
    height: calc(var(--page-size) * (2.25rem + var(--wa-space-xs)));
  }
</style>
