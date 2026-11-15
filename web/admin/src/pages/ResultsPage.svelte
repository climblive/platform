<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import "@awesome.me/webawesome/dist/components/breadcrumb-item/breadcrumb-item.js";
  import "@awesome.me/webawesome/dist/components/breadcrumb/breadcrumb.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { getCompetitionQuery } from "@climblive/lib/queries";
  import { navigate } from "svelte-routing";
  import ResultsList from "./ResultsList.svelte";

  interface Props {
    competitionId: string;
  }

  const { competitionId }: Props = $props();

  const competitionQuery = $derived(getCompetitionQuery(competitionId));
  const competition = $derived(competitionQuery.data);
</script>

{#if competition === undefined}
  <Loader />
{:else}
  <wa-breadcrumb>
    <wa-breadcrumb-item
      onclick={() =>
        navigate(
          `/admin/organizers/${competition.ownership.organizerId}/competitions`,
        )}><wa-icon name="home"></wa-icon></wa-breadcrumb-item
    >
    <wa-breadcrumb-item
      onclick={() => navigate(`/admin/competitions/${competitionId}`)}
      >{competition.name}</wa-breadcrumb-item
    >
    <wa-breadcrumb-item>Results</wa-breadcrumb-item>
  </wa-breadcrumb>

  <h1>Results</h1>

  <ResultsList {competitionId} />
{/if}

<style>
  wa-breadcrumb {
    margin-block-end: var(--wa-space-m);
    display: block;
  }
</style>
