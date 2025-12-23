<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import "@awesome.me/webawesome/dist/components/breadcrumb-item/breadcrumb-item.js";
  import "@awesome.me/webawesome/dist/components/breadcrumb/breadcrumb.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import type WaSwitch from "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import type { Contender } from "@climblive/lib/models";
  import {
    getContendersByContestQuery,
    getContestQuery,
  } from "@climblive/lib/queries";
  import { format } from "date-fns";
  import { Link, navigate } from "svelte-routing";

  interface Props {
    contestId: number;
  }

  const { contestId }: Props = $props();

  const contestQuery = $derived(getContestQuery(contestId));
  const contendersQuery = $derived(getContendersByContestQuery(contestId));

  const contest = $derived(contestQuery.data);
  const contenders = $derived(contendersQuery.data);

  let showUnusedOnly = $state(false);

  const filteredContenders = $derived.by(() => {
    if (!contenders) {
      return undefined;
    }

    if (showUnusedOnly) {
      return contenders.filter(({ entered }) => entered === undefined);
    }

    return contenders;
  });

  const handleToggleUnusedOnly = (event: InputEvent) => {
    showUnusedOnly = (event.target as WaSwitch).checked;
  };

  const columns: ColumnDefinition<Contender>[] = [
    {
      label: "Registration code",
      mobile: true,
      render: renderRegistrationCode,
      width: "1fr",
    },
    {
      label: "Name",
      mobile: true,
      render: renderName,
      width: "max-content",
    },
    {
      label: "Used",
      mobile: true,
      render: renderUsed,
      width: "max-content",
    },
    {
      label: "Entered",
      mobile: false,
      render: renderTimestamp,
      width: "max-content",
    },
  ];
</script>

{#snippet renderRegistrationCode({ id, registrationCode }: Contender)}
  <Link to={`/admin/contenders/${id}`}>{registrationCode}</Link>
{/snippet}

{#snippet renderName({ name }: Contender)}
  {name ?? "-"}
{/snippet}

{#snippet renderUsed({ entered }: Contender)}
  {#if entered}
    <wa-icon name="check"></wa-icon>
  {:else}
    -
  {/if}
{/snippet}

{#snippet renderTimestamp({ entered }: Contender)}
  {entered ? format(entered, "yyyy-MM-dd HH:mm") : "-"}
{/snippet}

{#if contest}
  <wa-breadcrumb>
    <wa-breadcrumb-item
      onclick={() =>
        navigate(`/admin/organizers/${contest.ownership.organizerId}`)}
      ><wa-icon name="home"></wa-icon></wa-breadcrumb-item
    >
    <wa-breadcrumb-item onclick={() => navigate(`/admin/contests/${contestId}`)}
      >{contest.name}</wa-breadcrumb-item
    >
    <wa-breadcrumb-item>Tickets</wa-breadcrumb-item>
  </wa-breadcrumb>

  <h1>Tickets</h1>

  <wa-switch
    size="small"
    checked={showUnusedOnly}
    onchange={handleToggleUnusedOnly}>Show unused only</wa-switch
  >

  {#if filteredContenders === undefined}
    <Loader />
  {:else if filteredContenders.length > 0}
    <Table {columns} data={filteredContenders} getId={({ id }) => id}></Table>
  {/if}
{/if}

<style>
  wa-breadcrumb {
    margin-block-end: var(--wa-space-m);
    display: block;
  }

  wa-switch {
    margin-block-end: var(--wa-space-m);
  }
</style>
