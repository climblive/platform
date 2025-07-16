<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/qr-code/qr-code.js";
  import { Table, TableCell, TableRow } from "@climblive/lib/components";
  import type { CreateContendersArguments } from "@climblive/lib/models";
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
</script>

{#if contenders}
  <p>
    You have {remainingCodes} codes remaining out of your maximum allotted 500.
  </p>
  <section>
    {#each increments as increment (increment)}
      <wa-button
        size="small"
        type="button"
        variant="primary"
        loading={$createContenders.isPending}
        disabled={!remainingCodes || remainingCodes < increment}
        onclick={() => addContenders({ number: increment })}
      >
        <wa-icon slot="slot" name="plus-lg"></wa-icon>
        Add {increment} code{#if increment != 1}s{/if}
      </wa-button>
    {/each}
  </section>

  <Table columns={["Code", "Name", "Placement", "Score"]}>
    {#each contenders as contender (contender.id)}
      <TableRow>
        <TableCell>
          <a href={`/${contender.registrationCode}`}>
            {contender.registrationCode}
          </a>
        </TableCell>
        <TableCell>{contender.name}</TableCell>
        <TableCell>{contender.score?.placement}</TableCell>
        <TableCell>{contender.score?.score}</TableCell>
      </TableRow>
    {/each}
  </Table>
{/if}

<style>
  section {
    display: flex;
    gap: var(--wa-space-xs);
  }
</style>
