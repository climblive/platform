<script lang="ts">
  import { Table, TableCell, TableRow } from "@climblive/lib/components";
  import type { CreateContendersArguments } from "@climblive/lib/models";
  import {
    createContendersMutation,
    getContendersByContestQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import "@shoelace-style/shoelace/dist/components/icon/icon.js";
  import "@shoelace-style/shoelace/dist/components/qr-code/qr-code.js";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const contendersQuery = getContendersByContestQuery(contestId);
  const createContenders = createContendersMutation(contestId);

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
      <sl-button
        size="small"
        type="button"
        variant="primary"
        loading={$createContenders.isPending}
        disabled={!remainingCodes || remainingCodes < increment}
        onclick={() => addContenders({ number: increment })}
      >
        <sl-icon slot="prefix" name="plus-lg"></sl-icon>
        Add {increment} code{#if increment != 1}s{/if}
      </sl-button>
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
    gap: var(--sl-spacing-x-small);
  }
</style>
