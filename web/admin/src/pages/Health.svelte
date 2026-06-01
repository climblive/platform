<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import RelativeTime from "@/components/RelativeTime.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import {
    EmptyState,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import type { Contest, ServiceStatus } from "@climblive/lib/models";
  import {
    getContestsQuery,
    getHealthQuery,
    getRunningScoreEnginesQuery,
    getVersionQuery,
    stopScoreEngineMutation,
    type RunningScoreEngine,
  } from "@climblive/lib/queries";
  import { Link } from "svelte-routing";

  const columns: ColumnDefinition<ServiceStatus>[] = [
    {
      mobile: true,
      render: renderStatus,
      width: "max-content",
    },
    {
      label: "Service",
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

  type RunningScoreEngineRow = RunningScoreEngine & {
    contest: Contest | undefined;
  };

  const runningScoreEngineColumns: ColumnDefinition<RunningScoreEngineRow>[] = [
    {
      label: "Engine ID",
      mobile: true,
      render: renderScoreEngineInstanceId,
      width: "1fr",
    },
    {
      label: "Contest",
      mobile: true,
      render: renderScoreEngineContest,
      width: "1fr",
    },
    {
      label: "Actions",
      mobile: true,
      render: renderScoreEngineActions,
      align: "right",
      width: "max-content",
    },
  ];

  const healthQuery = $derived(getHealthQuery());
  const health = $derived(healthQuery.data);
  const versionQuery = $derived(getVersionQuery());
  const version = $derived(versionQuery.data);
  const runningScoreEnginesQuery = $derived(getRunningScoreEnginesQuery());
  const runningScoreEngines = $derived(runningScoreEnginesQuery.data);
  const contestIds = $derived(
    runningScoreEngines?.map(({ contestId }) => contestId),
  );
  const contestsQuery = $derived(getContestsQuery(contestIds));
  const contests = $derived(contestsQuery.data);

  const runningScoreEngineRows = $derived.by(() => {
    const rows = (runningScoreEngines ?? []).map((engine) => ({
      ...engine,
      contest: contests?.find(({ id }) => id === engine.contestId),
    }));

    rows.sort((left, right) => left.contestId - right.contestId);
    return rows;
  });
  const allHealthy = $derived(health?.every(({ healthy }) => healthy));
  const stopScoreEngine = stopScoreEngineMutation();

  let confirmStopEngineId = $state<string | undefined>(undefined);

  const handleStopScoreEngine = (instanceId: string) => {
    if (confirmStopEngineId !== instanceId) {
      confirmStopEngineId = instanceId;
      return;
    }

    stopScoreEngine.mutate(instanceId, {
      onSettled: () => {
        confirmStopEngineId = undefined;
      },
    });
  };

  $effect(() => {
    if (
      confirmStopEngineId &&
      !runningScoreEngines?.some(
        ({ instanceId }) => instanceId === confirmStopEngineId,
      )
    ) {
      confirmStopEngineId = undefined;
    }
  });
</script>

{#snippet renderStatus({ healthy }: ServiceStatus)}
  {#if healthy}
    <wa-icon name="heart-circle-check" class="healthy"></wa-icon>
  {:else}
    <wa-icon name="disease" class="unhealthy"></wa-icon>
  {/if}
{/snippet}

{#snippet renderName({ name }: ServiceStatus)}
  {name}
{/snippet}

{#snippet renderLastSeen({ checkedAt }: ServiceStatus)}
  <RelativeTime time={checkedAt} />
{/snippet}

{#snippet renderScoreEngineInstanceId({ instanceId }: RunningScoreEngineRow)}
  {instanceId}
{/snippet}

{#snippet renderScoreEngineContest({
  contest,
  contestId,
}: RunningScoreEngineRow)}
  <Link to={`/admin/contests/${contestId}`}>
    {contest?.name ?? `Contest ${contestId}`}
  </Link>
{/snippet}

{#snippet renderScoreEngineActions({ instanceId }: RunningScoreEngineRow)}
  <wa-button
    size="s"
    appearance="outlined"
    variant={confirmStopEngineId === instanceId ? "danger" : "neutral"}
    loading={stopScoreEngine.isPending}
    onclick={() => handleStopScoreEngine(instanceId)}
  >
    {#if confirmStopEngineId === instanceId}
      Confirm stop
    {:else}
      Stop
    {/if}
    <wa-icon name="stop" slot="start"></wa-icon>
  </wa-button>
{/snippet}

<div class="title">
  <h1>System health</h1>
  {#if version !== undefined}
    <wa-badge pill variant="neutral">{version}</wa-badge>
  {/if}
</div>

{#if health === undefined}
  <Loader />
{:else}
  {#if allHealthy}
    <wa-callout variant="success">
      <wa-icon slot="icon" name="heart-circle-check"></wa-icon>
      All services are up and running.
    </wa-callout>
  {:else}
    <wa-callout variant="danger">
      <wa-icon slot="icon" name="disease"></wa-icon>
      One or more services are down.
    </wa-callout>
  {/if}
  <Table {columns} data={health} getId={({ name }) => name}></Table>

  <h2>Running score engines</h2>
  {#if contests === undefined || runningScoreEngines === undefined}
    <Loader />
  {:else if runningScoreEngineRows.length === 0}
    <EmptyState
      title="No score engines are currently running"
      description="Running score engines will appear here automatically."
    ></EmptyState>
  {:else}
    <Table
      columns={runningScoreEngineColumns}
      data={runningScoreEngineRows}
      getId={({ instanceId }) => instanceId}
    ></Table>
  {/if}
{/if}

<style>
  .title {
    display: flex;
    align-items: start;
    gap: var(--wa-space-m);
    margin-block: var(--wa-space-l);

    & h1 {
      margin-block: 0;
    }

    & wa-badge {
      font-size: var(--wa-font-size-3xs);
    }
  }

  .healthy {
    color: var(--wa-color-success);
  }

  .unhealthy {
    color: var(--wa-color-danger);
  }

  wa-callout {
    margin-block-end: var(--wa-space-m);
  }

  h2 {
    margin: var(--wa-space-xl) 0 var(--wa-space-m);
  }
</style>
