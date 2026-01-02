<script lang="ts">
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import type WaInput from "@awesome.me/webawesome/dist/components/input/input.js";
  import { value } from "@climblive/lib/forms";
  import type { CreateContendersArguments } from "@climblive/lib/models";
  import {
    createContendersMutation,
    getContendersByContestQuery,
    getContestQuery,
    patchContestMutation,
  } from "@climblive/lib/queries";
  import { toastError, toastSuccess } from "@climblive/lib/utils";
  import { Link } from "svelte-routing";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let dialog: WaDialog | undefined = $state();
  let unlockDialog: WaDialog | undefined = $state();
  let numberInput: WaInput | undefined = $state();

  let newTicketsAvailableForPrint = $state(false);

  const contestQuery = $derived(getContestQuery(contestId));
  const contendersQuery = $derived(getContendersByContestQuery(contestId));
  const createContenders = $derived(createContendersMutation(contestId));
  const patchContest = $derived(patchContestMutation(contestId));

  let contest = $derived(contestQuery.data);
  let contenders = $derived(contendersQuery.data);

  let maxTickets = $derived(contest?.evaluationMode ? 10 : 500);

  let remainingCodes = $derived(
    contenders === undefined ? undefined : maxTickets - contenders.length,
  );

  let registeredContenders = $derived.by(() => {
    if (!contenders) {
      return undefined;
    }

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

  const handleOpenUnlockDialog = () => {
    if (unlockDialog) {
      unlockDialog.open = true;
    }
  };

  const closeUnlockDialog = () => {
    if (unlockDialog) {
      unlockDialog.open = false;
    }
  };

  const handleCreate = () => {
    if (numberInput) {
      const args: CreateContendersArguments = {
        number: Number(numberInput.value),
      };

      createContenders.mutate(args, {
        onSuccess: () => {
          newTicketsAvailableForPrint = true;
          closeDialog();
        },
        onError: () => toastError("Failed to create tickets."),
      });
    }
  };

  const handleUnlockEvaluationMode = () => {
    patchContest.mutate(
      { evaluationMode: false },
      {
        onSuccess: () => {
          toastSuccess(
            "Evaluation mode unlocked. You can now create up to 500 tickets.",
          );
          closeUnlockDialog();
        },
        onError: () => toastError("Failed to unlock evaluation mode."),
      },
    );
  };
</script>

<p class="copy">
  Tickets hold unique registration codes, granting contenders access to your
  contest. These tickets may be printed on paper and distributed to the
  contenders on site.
  {#if contenders && contenders.length > 0}
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
      {@attach value(Math.min(100, remainingCodes ?? 0))}
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
    variant="neutral"
    appearance="accent"
    loading={createContenders.isPending}
    onclick={handleCreate}
    type="submit"
  >
    Create
  </wa-button>
</wa-dialog>

<wa-dialog bind:this={unlockDialog} label="Unlock full capacity">
  <div class="dialog-content">
    <wa-callout variant="warning">
      <wa-icon slot="icon" name="triangle-exclamation"></wa-icon>
      Help keep things neat and tidy by only unlocking this for real contests!
    </wa-callout>

    <p>
      Evaluation mode limits contests to 10 contenders for testing purposes. By
      unlocking the full capacity you'll be able to host up to 500 contenders.
    </p>
  </div>

  <wa-button slot="footer" appearance="plain" onclick={closeUnlockDialog}
    >Cancel</wa-button
  >
  <wa-button
    slot="footer"
    size="small"
    variant="success"
    appearance="accent"
    loading={patchContest.isPending}
    onclick={handleUnlockEvaluationMode}
  >
    <wa-icon slot="start" name="lock-open"></wa-icon>
    Unlock
  </wa-button>
</wa-dialog>

<div class="actions">
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
  {#if contest?.evaluationMode}
    <wa-button
      size="small"
      variant="success"
      appearance="accent"
      onclick={handleOpenUnlockDialog}
    >
      <wa-icon slot="start" name="lock-open"></wa-icon>
      Unlock evaluation mode
    </wa-button>
  {/if}
  {#if contenders && contenders.length > 0}
    <Link to={`/admin/contests/${contestId}/tickets`}>
      <wa-button appearance="outlined" size="small"
        >View all tickets
        <wa-icon name="list" slot="start"></wa-icon>
      </wa-button>
    </Link>
  {/if}
  <a href={`/admin/contests/${contestId}/tickets/print`} target="_blank">
    <wa-button
      appearance="outlined"
      size="small"
      loading={contendersQuery.isLoading}
      disabled={!contenders || contenders.length === 0}
      >Print tickets
      <wa-icon name="print" slot="start"></wa-icon>
      {#if contenders && contenders.length > 0}
        <wa-badge
          variant="neutral"
          attention={newTicketsAvailableForPrint ? "pulse" : undefined}
          pill>{contenders.length}</wa-badge
        >
      {/if}
    </wa-button>
  </a>
</div>

<p>
  {#if remainingCodes === maxTickets}
    You may create up to {maxTickets} tickets.
  {:else}
    You may create {remainingCodes} more tickets.
  {/if}
</p>

<style>
  .dialog-content {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }

  .actions {
    display: flex;
    gap: var(--wa-space-xs);
    flex-wrap: wrap;
  }

  .copy {
    color: var(--wa-color-text-quiet);
  }
</style>
