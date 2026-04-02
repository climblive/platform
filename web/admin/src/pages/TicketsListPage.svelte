<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import "@awesome.me/webawesome/dist/components/breadcrumb-item/breadcrumb-item.js";
  import "@awesome.me/webawesome/dist/components/breadcrumb/breadcrumb.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import type WaCheckbox from "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import type WaNumberInput from "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import type WaSwitch from "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { value } from "@climblive/lib/forms";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import type {
    Contender,
    CreateContendersArguments,
  } from "@climblive/lib/models";
  import {
    createContendersMutation,
    getContendersByContestQuery,
    getContestQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { Link, navigate } from "svelte-routing";

  const maxTickets = 500;

  interface Props {
    contestId: number;
  }

  const { contestId }: Props = $props();

  const contestQuery = $derived(getContestQuery(contestId));
  const contendersQuery = $derived(getContendersByContestQuery(contestId));
  const createContenders = createContendersMutation(contestId);

  const contest = $derived(contestQuery.data);
  const contenders = $derived(contendersQuery.data);

  let dialog: WaDialog | undefined = $state();
  let printDialog: WaDialog | undefined = $state();
  let numberInput: WaNumberInput | undefined = $state();

  const remainingCodes = $derived(
    contenders === undefined ? undefined : maxTickets - contenders.length,
  );

  let showUnusedOnly = $state(false);
  let selectionStartId: number | undefined = $state(undefined);
  let selectionEndId: number | undefined = $state(undefined);

  const filteredContenders = $derived.by(() => {
    if (!contenders) {
      return undefined;
    }

    const sorted = [...contenders].sort((a, b) => a.id - b.id);

    if (showUnusedOnly) {
      return sorted.filter(({ entered }) => entered === undefined);
    }

    return sorted;
  });

  const selectedCount = $derived.by(() => {
    const from = selectionStartId;
    const to = selectionEndId;

    if (from === undefined || to === undefined || !filteredContenders) {
      return 0;
    }

    return filteredContenders.filter((c) => c.id >= from && c.id <= to).length;
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
      selectionStartId = filteredContenders?.find((c) => c.id > id)?.id;
    } else if (id === selectionEndId) {
      selectionEndId = filteredContenders?.findLast((c) => c.id < id)?.id;
    }
  };

  const allSelected = $derived(
    filteredContenders !== undefined &&
      filteredContenders.length > 0 &&
      selectedCount === filteredContenders.length,
  );

  const handleToggleSelectAll = (event: InputEvent) => {
    const checkbox = event.target as WaCheckbox;

    if (
      checkbox.checked &&
      filteredContenders &&
      filteredContenders.length > 0
    ) {
      selectionStartId = filteredContenders[0].id;
      selectionEndId = filteredContenders[filteredContenders.length - 1].id;
    } else {
      selectionStartId = undefined;
      selectionEndId = undefined;
    }
  };

  const printUrl = $derived.by(() => {
    if (selectionStartId !== undefined && selectionEndId !== undefined) {
      return `/admin/contests/${contestId}/tickets/print?from=${selectionStartId}&to=${selectionEndId}`;
    }

    return undefined;
  });

  const handleOpenCreateDialog = async () => {
    if (dialog) {
      dialog.open = true;
    }
  };

  const closeDialog = () => {
    if (dialog) {
      dialog.open = false;
    }
  };

  const closePrintDialog = () => {
    if (printDialog) {
      printDialog.open = false;
    }
  };

  const handlePrintNewTickets = () => {
    closePrintDialog();

    if (printUrl) {
      window.open(printUrl, "_blank");
    }
  };

  const handleCreate = () => {
    if (numberInput) {
      const args: CreateContendersArguments = {
        number: Number(numberInput.value),
      };

      createContenders.mutate(args, {
        onSuccess: (newContenders) => {
          closeDialog();

          if (newContenders.length > 0) {
            const ids = newContenders.map((c) => c.id);
            selectionStartId = Math.min(...ids);
            selectionEndId = Math.max(...ids);

            if (printDialog) {
              printDialog.open = true;
            }
          }
        },
        onError: () => toastError("Failed to create tickets."),
      });
    }
  };

  const columns: ColumnDefinition<Contender>[] = [
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
      label: "Used",
      mobile: true,
      render: renderUsed,
      width: "minmax(5rem, 1fr)",
      align: "right",
    },
  ];
</script>

{#snippet renderSelectAll()}
  <wa-checkbox
    size="small"
    checked={allSelected}
    indeterminate={selectedCount > 0 && !allSelected}
    onchange={handleToggleSelectAll}
  ></wa-checkbox>
{/snippet}

{#snippet renderCheckbox(contender: Contender)}
  <wa-checkbox
    size="small"
    checked={isSelected(contender.id)}
    disabled={isLocked(contender.id)}
    onchange={(e: InputEvent) => handleCheckboxChange(contender.id, e)}
  ></wa-checkbox>
{/snippet}

{#snippet renderTicketNumber({ id }: Contender)}
  #{id.toString().padStart(6, "0")}
{/snippet}

{#snippet renderRegistrationCode({ id, registrationCode }: Contender)}
  <Link to={`/admin/contenders/${id}`}>
    <wa-icon name="qrcode"></wa-icon>
    <span class="regcode">{registrationCode}</span>
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
    <div class="selection-actions">
      <wa-button
        size="small"
        variant="neutral"
        appearance="accent"
        onclick={handleOpenCreateDialog}
        disabled={remainingCodes === undefined || remainingCodes === 0}
      >
        <wa-icon slot="start" name="plus"></wa-icon>
        Create tickets</wa-button
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

    <wa-switch
      size="small"
      checked={showUnusedOnly}
      onchange={handleToggleUnusedOnly}>Show unused only</wa-switch
    >
  </div>

  <wa-dialog bind:this={dialog} label="Create tickets">
    <div class="dialog-content">
      <wa-callout variant="neutral">
        <wa-icon slot="icon" name="circle-exclamation"></wa-icon>
        You have {remainingCodes} tickets remaining out of your maximum allotted
        {maxTickets}.
      </wa-callout>

      <wa-number-input
        bind:this={numberInput}
        name="number"
        {@attach value(Math.min(100, remainingCodes ?? 0))}
        min="1"
        max={remainingCodes}
        label="Number of tickets to create"
      ></wa-number-input>
    </div>

    <wa-button slot="footer" appearance="plain" onclick={closeDialog}
      >Cancel</wa-button
    >
    <wa-button
      slot="footer"
      size="small"
      variant="neutral"
      appearance="accent"
      loading={createContenders.isPending}
      onclick={handleCreate}
      type="submit"
    >
      Create
    </wa-button>
  </wa-dialog>

  <wa-dialog bind:this={printDialog} label="Print tickets">
    Do you want to print the newly created tickets?
    <wa-button slot="footer" appearance="plain" onclick={closePrintDialog}
      >Later</wa-button
    >
    <wa-button slot="footer" variant="neutral" onclick={handlePrintNewTickets}>
      <wa-icon slot="start" name="print"></wa-icon>
      Print
    </wa-button>
  </wa-dialog>

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

  .dialog-content {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }

  .regcode {
    font-family: monospace;
    font-size: var(--wa-font-size-m);
  }
</style>
