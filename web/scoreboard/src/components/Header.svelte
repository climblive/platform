<script lang="ts">
  import "@awesome.me/webawesome/dist/components/progress-bar/progress-bar.js";
  import { CompetitionStateProvider, Timer } from "@climblive/lib/components";
  import type { ScoreboardEntry } from "@climblive/lib/models";
  import { type Readable } from "svelte/store";

  interface Props {
    competitionId: number;
    name: string;
    categoryId: number;
    startTime: Date;
    endTime: Date;
    scoreboard: Readable<Map<number, ScoreboardEntry[]>>;
  }

  let {
    competitionId,
    name,
    categoryId,
    startTime,
    endTime,
    scoreboard,
  }: Props = $props();

  let categorySize = $derived(($scoreboard.get(categoryId) ?? []).length);
  let totalSize = $derived.by(() => {
    let total = 0;
    $scoreboard.forEach((entry) => (total += entry.length));
    return total;
  });
</script>

<CompetitionStateProvider {competitionId} {categoryId}>
  {#snippet children({ competitionState, progress })}
    <header>
      <div class="left">
        <div class="title">
          <small>Category</small>
          <h2>
            {name}
          </h2>
        </div>
        {#if competitionState === "NOT_STARTED"}
          <Timer endTime={startTime} label="Starting in" />
        {:else}
          <Timer {endTime} label="Time left" />
        {/if}
      </div>
      <div class="size">
        <strong>{categorySize}</strong>/{totalSize}
      </div>
      <wa-progress-bar value={progress}></wa-progress-bar>
    </header>
  {/snippet}
</CompetitionStateProvider>

<style>
  header {
    margin-bottom: var(--wa-space-s);
    background-color: var(--wa-color-surface-raised);
    border-radius: var(--wa-border-radius-m);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-surface-border);

    color: var(--wa-color-text-normal);

    display: grid;
    grid-template-columns: 1fr max-content;
    grid-template-rows: 1fr;

    .left {
      display: flex;
      flex-direction: column;
      gap: var(--wa-space-xs);
      align-items: start;
      min-width: 0;

      padding: var(--wa-space-s);
    }

    .title {
      width: 100%;
      min-width: 0;

      small {
        font-size: var(--wa-font-size-xs);
        color: var(--wa-color-text-quiet);
      }
    }

    & h2 {
      min-width: 0;
      margin: 0;
      font-weight: var(--wa-font-weight-bold);

      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .size {
      align-self: start;
      margin: var(--wa-space-s) var(--wa-space-s) 0 0;

      font-size: var(--wa-font-size-s);
      font-weight: var(--wa-font-weight-semibold);

      & strong {
        font-size: 1.5em;
      }
    }

    wa-progress-bar {
      --track-height: 0.5rem;

      grid-column: 1 / -1;
      padding: 0 var(--wa-space-s) var(--wa-space-s);
    }
  }
</style>
