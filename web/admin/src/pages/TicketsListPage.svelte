<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import "@awesome.me/webawesome/dist/components/breadcrumb-item/breadcrumb-item.js";
  import "@awesome.me/webawesome/dist/components/breadcrumb/breadcrumb.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import WaSwitch from "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import { checked } from "@climblive/lib/forms";
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

  let { contestId }: Props = $props();

  const contestQuery = $derived(getContestQuery(contestId));
  const contendersQuery = $derived(getContendersByContestQuery(contestId));

  const contest = $derived(contestQuery.data);
  const contenders = $derived(contendersQuery.data);

  let showUnusedOnly = $state(false);
  let unusedToggle: WaSwitch | undefined = $state();

  const filteredContenders = $derived.by(() => {
    if (!contenders) {
      return undefined;
    }

    if (!showUnusedOnly) {
      return contenders;
    }

    return contenders.filter((contender) => contender.entered === undefined);
  });

  const handleToggleUnusedOnly = () => {
    if (unusedToggle) {
      showUnusedOnly = unusedToggle.checked;
    }
  };

  const columns: ColumnDefinition<Contender>[] = [
    {
      label: "Registration Code",
      mobile: true,
      render: renderRegistrationCode,
    },
    {
      label: "Name",
      mobile: true,
      render: renderName,
    },
    {
      label: "Used",
      mobile: true,
      render: renderUsed,
    },
    {
      label: "Timestamp",
      mobile: false,
      render: renderTimestamp,
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
    <wa-icon name="minus"></wa-icon>
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

  <div class="controls">
    <wa-switch
      bind:this={unusedToggle}
      {@attach checked(showUnusedOnly)}
      onchange={handleToggleUnusedOnly}
      >Show unused only</wa-switch
    >
  </div>

  {#if filteredContenders === undefined}
    <Loader />
  {:else if filteredContenders.length > 0}
    <Table {columns} data={filteredContenders} getId={({ id }) => id}></Table>
  {:else}
    <p class="no-tickets">No tickets to display.</p>
  {/if}
{/if}

<style>
  wa-breadcrumb {
    margin-block-end: var(--wa-space-m);
    display: block;
  }

  .controls {
    margin-block-end: var(--wa-space-m);
  }

  .no-tickets {
    color: var(--wa-color-text-quiet);
  }
</style>
