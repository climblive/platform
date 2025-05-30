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
    margin-bottom: var(--sl-spacing-small);
    background-color: var(--sl-color-primary-600);
    border-radius: var(--sl-border-radius-medium);

    color: white;

    display: grid;
    grid-template-columns: 1fr max-content;
    grid-template-rows: 1fr;

    .left {
      display: flex;
      flex-direction: column;
      gap: var(--sl-spacing-x-small);
      align-items: start;

      padding: var(--sl-spacing-small);
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
      font-weight: var(--sl-font-weight-semibold);

      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .size {
      align-self: start;
      margin: var(--sl-spacing-small) var(--sl-spacing-small) 0 0;

      color: white;
      font-size: var(--sl-font-size-small);
      font-weight: var(--sl-font-weight-semibold);

      & strong {
        font-size: 1.5em;
      }
    }
  }
</style>
