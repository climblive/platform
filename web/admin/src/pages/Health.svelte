<script lang="ts">
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import Loader from "@/components/Loader.svelte";
  import RelativeTime from "@/components/RelativeTime.svelte";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import type { ServiceStatus } from "@climblive/lib/models";
  import { getHealthQuery, getVersionQuery } from "@climblive/lib/queries";

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

  const healthQuery = $derived(getHealthQuery());
  const health = $derived(healthQuery.data);
  const versionQuery = $derived(getVersionQuery());
  const version = $derived(versionQuery.data);

  const allHealthy = $derived(health?.every(({ healthy }) => healthy));
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
{/if}

<style>
  .healthy {
    color: var(--wa-color-success);
  }

  .unhealthy {
    color: var(--wa-color-danger);
  }

  wa-callout {
    margin-block-end: var(--wa-space-m);
  }

  .title {
    display: flex;
    align-items: center;
    gap: var(--wa-space-m);
  }
</style>
