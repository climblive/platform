<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
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
  import { toastError } from "@climblive/lib/utils";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const contendersQuery = $derived(getContendersByContestQuery(contestId));
  const createContenders = $derived(createContendersMutation(contestId));

  let contenders = $derived($contendersQuery.data);

  let remainingCodes = $derived(
    contenders === undefined ? undefined : 500 - contenders.length,
  );

  const increments = [1, 10, 100];

  const addContenders = async (args: CreateContendersArguments) => {
    $createContenders.mutate(args, {
      onError: () => toastError("Failed to create contenders."),
    });
  };

  const columns: ColumnDefinition<Contender>[] = [
    {
      label: "Code",
      mobile: true,
      render: renderRegistrationCode,
    },
    {
      label: "Name",
      mobile: true,
      render: renderName,
    },
    {
      label: "Score",
      mobile: true,
      render: renderScore,
    },
  ];
</script>

{#snippet renderRegistrationCode({ registrationCode }: Contender)}
  <a href={`/${registrationCode}`}>
    {registrationCode}
  </a>
{/snippet}

{#snippet renderName({ name }: Contender)}
  {name}
{/snippet}

{#snippet renderScore({ score }: Contender)}
  {#if score}
    {score.placement} ({score.score})
  {/if}
{/snippet}

{#if contenders}
  <p>
    You have {remainingCodes} codes remaining out of your maximum allotted 500.
  </p>
  <section>
    {#each increments as increment (increment)}
      <wa-button
        size="small"
        variant="brand"
        appearance="accent"
        loading={$createContenders.isPending}
        disabled={!remainingCodes || remainingCodes < increment}
        onclick={() => addContenders({ number: increment })}
      >
        <wa-icon slot="start" name="plus-lg"></wa-icon>
        Add {increment} code{#if increment != 1}s{/if}
      </wa-button>
    {/each}
  </section>

  <Table {columns} data={contenders} getId={({ id }) => id}></Table>
{/if}

<style>
  section {
    display: flex;
    gap: var(--wa-space-xs);
  }
</style>
