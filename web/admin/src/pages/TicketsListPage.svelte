<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import "@awesome.me/webawesome/dist/components/breadcrumb-item/breadcrumb-item.js";
  import "@awesome.me/webawesome/dist/components/breadcrumb/breadcrumb.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import type WaCheckbox from "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import type WaSwitch from "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import type { Contender } from "@climblive/lib/models";
  import {
    getContendersByContestQuery,
    getContestQuery,
  } from "@climblive/lib/queries";
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
  let selectionStartId: number | undefined = $state(undefined);
  let selectionEndId: number | undefined = $state(undefined);

  const filteredContenders = $derived.by(() => {
    if (!contenders) {
      return undefined;
    }

    if (showUnusedOnly) {
      return contenders.filter(({ entered }) => entered === undefined);
    }

    return contenders;
  });

  const selectedCount = $derived.by(() => {
    const startId = selectionStartId;
    const endId = selectionEndId;

    if (startId === undefined || endId === undefined || !filteredContenders) {
      return 0;
    }

    return filteredContenders.filter((c) => c.id >= startId && c.id <= endId)
      .length;
  });

  const isSelected = (id: number): boolean => {
    return (
      selectionStartId !== undefined &&
      selectionEndId !== undefined &&
      id >= selectionStartId &&
      id <= selectionEndId
    );
  };

  const handleToggleUnusedOnly = (event: InputEvent) => {
    showUnusedOnly = (event.target as WaSwitch).checked;
    selectionStartId = undefined;
    selectionEndId = undefined;
  };

  const handleCheckboxChange = (id: number, event: InputEvent) => {
    const checkbox = event.target as WaCheckbox;
    const wantsChecked = checkbox.checked;

    if (wantsChecked) {
      if (selectionStartId === undefined || selectionEndId === undefined) {
        selectionStartId = id;
        selectionEndId = id;
      } else {
        selectionStartId = Math.min(selectionStartId, id);
        selectionEndId = Math.max(selectionEndId, id);
      }
    } else {
      if (selectionStartId === selectionEndId) {
        selectionStartId = undefined;
        selectionEndId = undefined;
      } else if (id === selectionStartId) {
        const startId = selectionStartId;
        const next = filteredContenders?.find((c) => c.id > startId);
        if (next) {
          selectionStartId = next.id;
        }
      } else if (id === selectionEndId) {
        const endId = selectionEndId;
        const prev = [...(filteredContenders ?? [])]
          .reverse()
          .find((c) => c.id < endId);
        if (prev) {
          selectionEndId = prev.id;
        }
      } else {
        checkbox.checked = true;
      }
    }
  };

  const handleSelectAll = () => {
    if (filteredContenders && filteredContenders.length > 0) {
      selectionStartId = filteredContenders[0].id;
      selectionEndId = filteredContenders[filteredContenders.length - 1].id;
    }
  };

  const handleDeselectAll = () => {
    selectionStartId = undefined;
    selectionEndId = undefined;
  };

  const printUrl = $derived.by(() => {
    if (selectionStartId !== undefined && selectionEndId !== undefined) {
      return `/admin/contests/${contestId}/tickets/print?from=${selectionStartId}&to=${selectionEndId}`;
    }

    return `/admin/contests/${contestId}/tickets/print`;
  });

  const columns: ColumnDefinition<Contender>[] = [
    {
      label: "",
      mobile: true,
      render: renderCheckbox,
      width: "max-content",
    },
    {
      label: "№",
      mobile: true,
      render: renderTicketNumber,
      width: "max-content",
    },
    {
      label: "Code",
      mobile: true,
      render: renderRegistrationCode,
      width: "max-content",
    },
    {
      label: "Used",
      mobile: true,
      render: renderUsed,
      width: "1fr",
      align: "right",
    },
  ];
</script>

{#snippet renderCheckbox(contender: Contender)}
  <wa-checkbox
    size="small"
    checked={isSelected(contender.id)}
    onchange={(e: InputEvent) => handleCheckboxChange(contender.id, e)}
  ></wa-checkbox>
{/snippet}

{#snippet renderTicketNumber({ id }: Contender)}
  #{id.toString().padStart(6, "0")}
{/snippet}

{#snippet renderRegistrationCode({ id, registrationCode }: Contender)}
  <Link to={`/admin/contenders/${id}`}>
    <wa-icon name="qrcode"></wa-icon>
    {registrationCode}
  </Link>
{/snippet}

{#snippet renderUsed({ entered }: Contender)}
  {#if entered}
    <wa-icon name="check"></wa-icon>
  {:else}
    -
  {/if}
{/snippet}

{#if contest}
  <wa-breadcrumb>
    <wa-breadcrumb-item
      onclick={() =>
        navigate(`/admin/organizers/${contest.ownership.organizerId}/contests`)}
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
      size="small"
      checked={showUnusedOnly}
      onchange={handleToggleUnusedOnly}>Show unused only</wa-switch
    >

    <div class="selection-actions">
      <wa-button
        size="small"
        appearance="outlined"
        onclick={handleSelectAll}
        disabled={!filteredContenders || filteredContenders.length === 0}
        >Select all</wa-button
      >
      <wa-button
        size="small"
        appearance="outlined"
        onclick={handleDeselectAll}
        disabled={selectedCount === 0}>Deselect</wa-button
      >
      <a href={printUrl} target="_blank">
        <wa-button
          size="small"
          appearance="outlined"
          disabled={selectedCount === 0}
        >
          <wa-icon name="print" slot="start"></wa-icon>
          Print selected
          {#if selectedCount > 0}
            <wa-badge variant="neutral" pill>{selectedCount}</wa-badge>
          {/if}
        </wa-button>
      </a>
    </div>
  </div>

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

  .controls {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: var(--wa-space-s);
    margin-block-end: var(--wa-space-m);
  }

  .selection-actions {
    display: flex;
    gap: var(--wa-space-xs);
    align-items: center;
  }
</style>
