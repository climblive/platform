<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import RelativeTime from "@/components/RelativeTime.svelte";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import type { RunnerStatus } from "@climblive/lib/models";
  import { getHealthQuery } from "@climblive/lib/queries";

  const healthQuery = $derived(getHealthQuery());
  const health = $derived(healthQuery.data);
</script>

{#snippet runnerRow(name: string, status: RunnerStatus)}
  <div class="runner">
    <span class="runner-name">{name}</span>
    {#if status.healthy}
      <wa-icon name="circle-check" class="healthy"></wa-icon>
    {:else}
      <wa-icon name="circle-xmark" class="unhealthy"></wa-icon>
    {/if}
    <span class="checked-at">
      {#if status.checkedAt}
        <RelativeTime time={status.checkedAt} />
      {:else}
        Never
      {/if}
    </span>
  </div>
{/snippet}

<h1>System Health</h1>

{#if health === undefined}
  <Loader />
{:else}
  <div class="runners">
    {@render runnerRow("Score Engine Manager", health.scoreEngineManager)}
    {@render runnerRow("Score Keeper", health.scoreKeeper)}
    {@render runnerRow("Scrubber", health.scrubber)}
  </div>
{/if}

<style>
  .runners {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }

  .runner {
    display: grid;
    grid-template-columns: 1fr max-content max-content;
    align-items: center;
    gap: var(--wa-space-m);
    padding: var(--wa-space-s) var(--wa-space-m);
    border: 1px solid var(--wa-color-surface-border);
    border-radius: var(--wa-border-radius-m);
    background-color: var(--wa-color-surface-raised);
  }

  .runner-name {
    font-weight: var(--wa-font-weight-semibold);
  }

  .checked-at {
    color: var(--wa-color-text-quiet);
    font-size: var(--wa-font-size-s);
    min-width: 10rem;
    text-align: right;
  }

  :global(.healthy) {
    color: var(--wa-color-success);
  }

  :global(.unhealthy) {
    color: var(--wa-color-danger);
  }
</style>
