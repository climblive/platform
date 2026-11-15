<script lang="ts">
  import CreateTicketsDialog from "@/components/CreateTicketsDialog.svelte";
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
  import {
    RegistrationCode,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import type { Competitor } from "@climblive/lib/models";
  import {
    getCompetitorsByCompetitionQuery,
    getCompetitionQuery,
  } from "@climblive/lib/queries";
  import { Link, navigate } from "svelte-routing";

  const maxTickets = 500;

  interface Props {
    competitionId: number;
  }

  const { competitionId }: Props = $props();

  const competitionQuery = $derived(getCompetitionQuery(competitionId));
  const competitorsQuery = $derived(
    getCompetitorsByCompetitionQuery(competitionId),
  );

  const competition = $derived(competitionQuery.data);
  const competitors = $derived(competitorsQuery.data);

  let createTicketsDialog: CreateTicketsDialog | undefined = $state();

  const remainingCodes = $derived(
    competitors === undefined ? undefined : maxTickets - competitors.length,
  );

  let showUnusedOnly = $state(false);
  let selectionStartId: number | undefined = $state(undefined);
  let selectionEndId: number | undefined = $state(undefined);

  const filteredCompetitors = $derived.by(() => {
    if (!competitors) {
      return undefined;
    }

    const sorted = [...competitors].sort((a, b) => a.id - b.id);

    if (showUnusedOnly) {
      return sorted.filter(({ entered }) => entered === undefined);
    }

    return sorted;
  });

  const selectedCount = $derived.by(() => {
    const from = selectionStartId;
    const to = selectionEndId;

    if (from === undefined || to === undefined || !filteredCompetitors) {
      return 0;
    }

    return filteredCompetitors.filter(({ id }) => isSelected(id)).length;
  });

  const isSelected = (id: number): boolean => {
    return (
      selectionStartId !== undefined &&
      selectionEndId !== undefined &&
      id >= selectionStartId &&
      id <= selectionEndId
    );
  };

  const isLocked = (id: number): boolean => {
    return isSelected(id) && id !== selectionStartId && id !== selectionEndId;
  };

  const handleToggleUnusedOnly = (event: InputEvent) => {
    showUnusedOnly = (event.target as WaSwitch).checked;
    selectionStartId = undefined;
    selectionEndId = undefined;
  };

  const handleCheckboxChange = (id: number, event: InputEvent) => {
    const checkbox = event.target as WaCheckbox;

    if (checkbox.checked) {
      selectionStartId = Math.min(selectionStartId ?? id, id);
      selectionEndId = Math.max(selectionEndId ?? id, id);
    } else if (id === selectionStartId && id === selectionEndId) {
      selectionStartId = undefined;
      selectionEndId = undefined;
    } else if (id === selectionStartId) {
      selectionStartId = filteredCompetitors?.find((c) => c.id > id)?.id;
    } else if (id === selectionEndId) {
      selectionEndId = filteredCompetitors?.findLast((c) => c.id < id)?.id;
    }
  };

  const allSelected = $derived(
    filteredCompetitors !== undefined &&
      filteredCompetitors.length > 0 &&
      selectedCount === filteredCompetitors.length,
  );

  const handleToggleSelectAll = (event: InputEvent) => {
    const checkbox = event.target as WaCheckbox;

    if (
      checkbox.checked &&
      filteredCompetitors &&
      filteredCompetitors.length > 0
    ) {
      selectionStartId = filteredCompetitors[0].id;
      selectionEndId = filteredCompetitors[filteredCompetitors.length - 1].id;
    } else {
      selectionStartId = undefined;
      selectionEndId = undefined;
    }
  };

  const printUrl = $derived.by(() => {
    if (selectionStartId === undefined || selectionEndId === undefined) {
      return undefined;
    }

    return `/admin/competitions/${competitionId}/tickets/print?from=${selectionStartId}&to=${selectionEndId}`;
  });

  const handleCreated = (newCompetitors: Competitor[]) => {
    const ids = newCompetitors.map(({ id }) => id);
    selectionStartId = Math.min(...ids);
    selectionEndId = Math.max(...ids);
  };

  const columns: ColumnDefinition<Competitor>[] = [
    {
      label: renderSelectAll,
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
      label: "Status",
      mobile: true,
      render: renderUsed,
      width: "1fr",
      align: "right",
    },
  ];
</script>

{#snippet renderSelectAll()}
  <wa-checkbox
    size="s"
    checked={allSelected}
    indeterminate={selectedCount > 0 && !allSelected}
    onchange={handleToggleSelectAll}
  ></wa-checkbox>
{/snippet}

{#snippet renderCheckbox(competitor: Competitor)}
  <wa-checkbox
    size="s"
    checked={isSelected(competitor.id)}
    disabled={isLocked(competitor.id)}
    onchange={(e: InputEvent) => handleCheckboxChange(competitor.id, e)}
  ></wa-checkbox>
{/snippet}

{#snippet renderTicketNumber({ id }: Competitor)}
  #{id.toString().padStart(6, "0")}
{/snippet}

{#snippet renderRegistrationCode({ id, registrationCode }: Competitor)}
  <Link to={`/admin/competitors/${id}`}>
    <wa-icon name="qrcode"></wa-icon>
    <span class="regcode"><RegistrationCode code={registrationCode} /></span>
  </Link>
{/snippet}

{#snippet renderUsed({ entered }: Competitor)}
  {#if entered}
    <wa-badge variant="success"
      ><wa-icon slot="start" name="user-check"></wa-icon>Registered</wa-badge
    >
  {:else}
    <wa-badge variant="neutral" appearance="filled">Unused</wa-badge>
  {/if}
{/snippet}

{#if competition}
  <wa-breadcrumb>
    <wa-breadcrumb-item
      onclick={() =>
        navigate(
          `/admin/organizers/${competition.ownership.organizerId}/competitions`,
        )}><wa-icon name="home"></wa-icon></wa-breadcrumb-item
    >
    <wa-breadcrumb-item
      onclick={() => navigate(`/admin/competitions/${competitionId}`)}
      >{competition.name}</wa-breadcrumb-item
    >
    <wa-breadcrumb-item>Tickets</wa-breadcrumb-item>
  </wa-breadcrumb>

  <h1>Tickets</h1>

  <p class="copy">
    Tickets contain registration codes that allow competitors to enter your
    competition. These tickets may be printed on paper and distributed on site.
  </p>

  <div class="controls">
    <div class="selection-actions">
      <wa-button
        size="s"
        variant="neutral"
        appearance="accent"
        onclick={() => createTicketsDialog?.open()}
        disabled={remainingCodes === undefined || remainingCodes === 0}
      >
        <wa-icon slot="start" name="plus"></wa-icon>
        Create tickets</wa-button
      >
      <a href={printUrl} target="_blank">
        <wa-button
          size="s"
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

    <wa-switch
      size="s"
      checked={showUnusedOnly}
      onchange={handleToggleUnusedOnly}>Show unused only</wa-switch
    >
  </div>

  <CreateTicketsDialog
    bind:this={createTicketsDialog}
    {competitionId}
    {remainingCodes}
    onCreated={handleCreated}
  />

  {#if filteredCompetitors === undefined}
    <Loader />
  {:else if filteredCompetitors.length > 0}
    <Table {columns} data={filteredCompetitors} getId={({ id }) => id}></Table>
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

  .regcode {
    font-family: monospace;
    font-size: var(--wa-font-size-m);
  }

  .copy {
    color: var(--wa-color-text-quiet);
    margin-bottom: var(--wa-space-m);
  }

  wa-badge {
    font-size: var(--wa-font-size-xs);
  }
</style>
