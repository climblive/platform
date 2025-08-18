<script lang="ts">
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import type WaInput from "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/qr-code/qr-code.js";
  import type { CreateContendersArguments } from "@climblive/lib/models";
  import {
    createContendersMutation,
    getContendersByContestQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";

  const maxTickets = 500;

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let dialog: WaDialog | undefined = $state();
  let numberInput: WaInput | undefined = $state();

  let newTicketsAvailableForPrint = $state(false);

  const contendersQuery = $derived(getContendersByContestQuery(contestId));
  const createContenders = $derived(createContendersMutation(contestId));

  let contenders = $derived($contendersQuery.data);

  let remainingCodes = $derived(
    contenders === undefined ? undefined : maxTickets - contenders.length,
  );

  let registeredContenders = $derived.by(() => {
    if (!contenders) return undefined;

    let count = 0;

    for (const contender of contenders) {
      if (contender.entered !== undefined) {
        count += 1;
      }
    }

    return count;
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

  const handleCreate = () => {
    if (numberInput) {
      const args: CreateContendersArguments = {
        number: Number(numberInput.value),
      };

      $createContenders.mutate(args, {
        onSuccess: () => {
          newTicketsAvailableForPrint = true;
          closeDialog();
        },
        onError: () => toastError("Failed to create contenders."),
      });
    }
  };
</script>

{#if contenders}
  <p class="copy">
    Tickets hold unique registration codes, granting contenders access to your
    contest. These tickets may be printed on paper and distributed to the
    contenders on site.
    {#if contenders.length > 0}
      Out of the {contenders.length}
      tickets that you have created, {registeredContenders} have already been used.
    {/if}
  </p>

  <wa-dialog bind:this={dialog} label="Create tickets">
    <div class="dialog-content">
      <wa-callout variant="neutral">
        <wa-icon slot="icon" name="circle-exclamation"></wa-icon>
        You have {remainingCodes} tickets remaining out of your maximum allotted
        {maxTickets}.
      </wa-callout>

      <wa-input
        bind:this={numberInput}
        name="number"
        type="number"
        value="100"
        min="1"
        max={remainingCodes}
        label="Number of tickets to create"
      ></wa-input>
    </div>

    <wa-button slot="footer" appearance="plain" onclick={closeDialog}
      >Cancel</wa-button
    >
    <wa-button
      slot="footer"
      size="small"
      variant="brand"
      appearance="accent"
      loading={$createContenders.isPending}
      onclick={handleCreate}
      type="submit"
    >
      Create
    </wa-button>
  </wa-dialog>

  <div class="actions">
    <wa-button
      size="small"
      variant="brand"
      appearance="accent"
      onclick={handleOpenCreateDialog}
    >
      <wa-icon slot="start" name="plus"></wa-icon>
      Create tickets</wa-button
    >
    {#if contenders.length > 0}
      <a href={`/admin/contests/${contestId}/tickets?print`} target="_blank">
        <wa-button appearance="outlined" size="small"
          >Print tickets
          <wa-icon name="print" slot="start"></wa-icon>
          <wa-badge
            variant="neutral"
            attention={newTicketsAvailableForPrint ? "pulse" : undefined}
            pill>{contenders.length}</wa-badge
          >
        </wa-button>
      </a>
    {/if}
  </div>

  <p>
    {#if remainingCodes === maxTickets}
      You may create up to {maxTickets} tickets.
    {:else}
      You may create {remainingCodes} more tickets.
    {/if}
  </p>
{/if}

<style>
  .dialog-content {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }

  .actions {
    display: flex;
    gap: var(--wa-space-xs);
  }

  .copy {
    color: var(--wa-color-text-quiet);
  }
</style>
