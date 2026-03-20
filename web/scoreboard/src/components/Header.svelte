<script lang="ts">
  import "@awesome.me/webawesome/dist/components/progress-bar/progress-bar.js";
  import { ContestStateProvider, Timer } from "@climblive/lib/components";
  import type { ScoreboardEntry } from "@climblive/lib/models";
  import { type Readable } from "svelte/store";

  interface Props {
    name: string;
    compClassId: number;
    startTime: Date;
    endTime: Date;
    scoreboard: Readable<Map<number, ScoreboardEntry[]>>;
  }

  let { name, compClassId, startTime, endTime, scoreboard }: Props = $props();

  let classSize = $derived(($scoreboard.get(compClassId) ?? []).length);
  let totalSize = $derived.by(() => {
    let total = 0;
    $scoreboard.forEach((entry) => (total += entry.length));
    return total;
  });
</script>

<ContestStateProvider {startTime} {endTime}>
  {#snippet children({ contestState, progress })}
    <header>
      <div class="left">
        <div class="title">
          <small>Class</small>
          <h2>
            {name}
          </h2>
        </div>
        {#if contestState === "NOT_STARTED"}
          <Timer endTime={startTime} label="Starting in" />
        {:else}
          <Timer {endTime} label="Time left" />
        {/if}
      </div>
      <div class="size">
        <strong>{classSize}</strong>/{totalSize}
      </div>
      <div class="progress">
        <wa-progress-bar value={progress}></wa-progress-bar>
      </div>
    </header>
  {/snippet}
</ContestStateProvider>

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

      padding: var(--wa-space-s);
    }

    .title {
      width: 100%;
      position: relative;

      small {
        font-size: var(--wa-font-size-xs);
        color: var(--wa-color-text-quiet);
      }
    }

    & h2 {
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

    .progress {
      grid-column: 1 / -1;
      padding: 0 var(--wa-space-s) var(--wa-space-s);
    }

    wa-progress-bar {
      --height: 2px;
    }
  }
</style>
