<script lang="ts">
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
  {#snippet children({ contestState })}
    <header>
      <div class="left">
        <div class="title">
          <h2>
            {name}
          </h2>
        </div>
        {#if contestState === "NOT_STARTED"}
          <Timer endTime={startTime} label="Time until start" />
        {:else}
          <Timer {endTime} label="Time remaining" />
        {/if}
      </div>
      {#if classSize > 0}
        <div class="size">
          <strong>{classSize}</strong>/{totalSize}
        </div>
      {/if}
    </header>
  {/snippet}
</ContestStateProvider>

<style>
  header {
    margin-bottom: var(--wa-space-s);
    background-color: var(--wa-color-brand-fill-normal);
    border-radius: var(--wa-panel-border-radius);
    border: var(--wa-panel-border-width) var(--wa-panel-border-style)
      var(--wa-color-brand-border-normal);

    color: var(--wa-color-brand-on-normal);

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
      height: 2rem;
    }

    & h2 {
      position: absolute;
      inset: 0;

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
  }
</style>
