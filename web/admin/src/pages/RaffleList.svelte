<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import {
    EmptyState,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
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

  let raffles = $derived(rafflesQuery.data);

  const handleCreateRaffle = () => {
    createRaffle.mutate(undefined, {
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

<section>
  {#if raffles === undefined}
    <Loader />
  {:else if raffles.length > 0}
    <wa-button variant="neutral" appearance="accent" onclick={handleCreateRaffle}
      >Start new raffle</wa-button
    >
    <Table {columns} data={raffles} getId={({ id }) => id}></Table>
  {:else}
    <EmptyState
      title="No raffles yet"
      description="Start a new raffle to randomly select winners from your contenders."
    >
      {#snippet actions()}
        <wa-button variant="neutral" appearance="accent" onclick={handleCreateRaffle}
          >Start new raffle</wa-button
        >
      {/snippet}
    </EmptyState>
  {/if}
</section>

<style>
  section {
    display: flex;
    flex-direction: column;
    align-items: start;
    gap: var(--wa-space-m);
  }
</style>
