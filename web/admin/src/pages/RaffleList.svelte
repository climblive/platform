<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
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

  const columns: ColumnDefinition<Raffle>[] = [
    {
      label: "Name",
      mobile: true,
      render: renderName,
    },
  ];
</script>

{#snippet renderName({ id }: Raffle)}
  <Link to={`/admin/raffles/${id}`}>Raffle {id}</Link>
{/snippet}

<wa-button variant="brand" appearance="accent" onclick={handleCreateRaffle}
  >Create</wa-button
>

<section>
  {#if raffles?.length}
    <Table {columns} data={raffles} getId={({ id }) => id}></Table>
  {/if}
</section>

<style>
  section {
    display: flex;
    gap: var(--wa-space-xs);
  }
</style>
