<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import type WaInput from "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/qr-code/qr-code.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import type {
    Contender,
    CreateContendersArguments,
  } from "@climblive/lib/models";
  import {
    createContendersMutation,
    getContendersByContestQuery,
  } from "@climblive/lib/queries";
  import { ordinalSuperscript, toastError } from "@climblive/lib/utils";
  import { Link } from "svelte-routing";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let dialog: WaDialog | undefined = $state();
  let numberInput: WaInput | undefined = $state();

  const contendersQuery = $derived(getContendersByContestQuery(contestId));
  const createContenders = $derived(createContendersMutation(contestId));

  let contenders = $derived($contendersQuery.data);

  let remainingCodes = $derived(
    contenders === undefined ? undefined : 500 - contenders.length,
  );

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

  const columns: ColumnDefinition<Contender>[] = [
    {
      label: "Code",
      mobile: false,
      render: renderRegistrationCode,
      width: "max-content",
    },
    {
      label: "Name",
      mobile: true,
      render: renderName,
      width: "3fr",
    },
    {
      label: "Score",
      mobile: true,
      render: renderScore,
      width: "max-content",
      align: "right",
    },
    {
      label: "Placement",
      mobile: true,
      render: renderPlacement,
      width: "max-content",
      align: "right",
    },
  ];

  const handleCreate = () => {
    if (numberInput) {
      const args: CreateContendersArguments = {
        number: Number(numberInput.value),
      };

      $createContenders.mutate(args, {
        onSuccess: closeDialog,
        onError: () => toastError("Failed to create contenders."),
      });
    }
  };
</script>

{#snippet renderRegistrationCode({ registrationCode }: Contender)}
  <a href={`/${registrationCode}`} target="blank">
    {registrationCode}
  </a>
{/snippet}

{#snippet renderName({ name }: Contender)}
  {name}
{/snippet}

{#snippet renderScore({ score }: Contender)}
  {#if score}
    {score.score}
  {/if}
{/snippet}

{#snippet renderPlacement({ score }: Contender)}
  {#if score}
    {score.placement}<sup>{ordinalSuperscript(score.placement)}</sup>
  {/if}
{/snippet}

{#if contenders}
  <p>
    You have {remainingCodes} codes remaining out of your maximum allotted 500.
  </p>
  <wa-dialog bind:this={dialog} label="Create contender tickets">
    <div class="dialog-content">
      <wa-callout variant="neutral">
        <wa-icon slot="icon" name="circle-exclamation"></wa-icon>
        You have {remainingCodes} codes remaining out of your maximum allotted 500.
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
    <Link to={`/admin/contests/${contestId}/tickets`}>
      <wa-button appearance="outlined" size="small"
        >Print tickets
        <wa-icon name="print" slot="start"></wa-icon>
      </wa-button>
    </Link>
  </div>

  <Table {columns} data={contenders} getId={({ id }) => id}></Table>
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
</style>
