<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import RelativeTime from "@/components/RelativeTime.svelte";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import type { RunnerStatus } from "@climblive/lib/models";
  import { getHealthQuery } from "@climblive/lib/queries";

  type RunnerRow = { name: string; status: RunnerStatus };

  const columns: ColumnDefinition<RunnerRow>[] = [
    {
      label: "Status",
      mobile: true,
      render: renderStatus,
      width: "max-content",
    },
    {
      label: "Runner",
      mobile: true,
      render: renderName,
      width: "1fr",
    },
    {
      label: "Last seen",
      mobile: true,
      render: renderLastSeen,
      align: "right",
      width: "max-content",
    },
  ];

  const healthQuery = $derived(getHealthQuery());
  const health = $derived(healthQuery.data);

  const rows = $derived<RunnerRow[]>(
    health === undefined
      ? []
      : [
          { name: "Score engine manager", status: health.scoreEngineManager },
          { name: "Score keeper", status: health.scoreKeeper },
          { name: "Scrubber", status: health.scrubber },
        ],
  );
</script>

{#snippet renderStatus({ status }: RunnerRow)}
  {#if status.healthy}
    <wa-icon name="circle-check" class="healthy"></wa-icon>
  {:else}
    <wa-icon name="circle-xmark" class="unhealthy"></wa-icon>
  {/if}
{/snippet}

{#snippet renderName({ name }: RunnerRow)}
  {name}
{/snippet}

{#snippet renderLastSeen({ status }: RunnerRow)}
  {#if status.checkedAt}
    <RelativeTime time={status.checkedAt} />
  {:else}
    Never
  {/if}
{/snippet}

<h1>System health</h1>

{#if health === undefined}
  <Loader />
{:else}
  <Table {columns} data={rows} getId={({ name }) => name}></Table>
{/if}

<style>
  :global(.healthy) {
    color: var(--wa-color-success);
  }

  :global(.unhealthy) {
    color: var(--wa-color-danger);
  }
</style>
