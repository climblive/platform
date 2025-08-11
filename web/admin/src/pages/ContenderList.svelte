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
  import { ordinalSuperscript, toastError } from "@climblive/lib/utils";

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
    },
    {
      label: "Placement",
      mobile: true,
      render: renderPlacement,
      width: "max-content",
    },
    {
      label: "Code",
      mobile: false,
      render: renderRegistrationCode,
      width: "max-content",
    },
  ];
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
