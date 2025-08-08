<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import { Table, TableCell, TableRow } from "@climblive/lib/components";
  import type { Raffle } from "@climblive/lib/models";
  import {
    createRaffleMutation,
    getRafflesQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { Link, navigate } from "svelte-routing";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const rafflesQuery = $derived(getRafflesQuery(contestId));
  const createRaffle = $derived(createRaffleMutation(contestId));

  let raffles = $derived($rafflesQuery.data);

  const handleCreateRaffle = () => {
    $createRaffle.mutate(undefined, {
      onSuccess: (raffle: Raffle) => navigate(`/admin/raffles/${raffle.id}`),
      onError: () => toastError("Failed to create raffle."),
    });
  };
</script>

<wa-button variant="brand" appearance="accent" onclick={handleCreateRaffle}
  >Create</wa-button
>

<section>
  <Table columns={["Name"]}>
    {#if raffles}
      {#each raffles as raffle (raffle.id)}
        <TableRow>
          <TableCell
            ><Link to={`/admin/raffles/${raffle.id}`}>Raffle {raffle.id}</Link
            ></TableCell
          >
        </TableRow>
      {/each}
    {/if}
  </Table>
</section>

<style>
  section {
    display: flex;
    gap: var(--wa-space-xs);
  }
</style>
