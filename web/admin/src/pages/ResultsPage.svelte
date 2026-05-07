<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import "@awesome.me/webawesome/dist/components/breadcrumb-item/breadcrumb-item.js";
  import "@awesome.me/webawesome/dist/components/breadcrumb/breadcrumb.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { getContestQuery } from "@climblive/lib/queries";
  import { navigate } from "svelte-routing";
  import ResultsList from "./ResultsList.svelte";

  interface Props {
    contestId: number;
  }

  const { contestId }: Props = $props();

  const contestQuery = $derived(getContestQuery(contestId));
  const contest = $derived(contestQuery.data);
</script>

{#if contest === undefined}
  <Loader />
{:else}
  <wa-breadcrumb>
    <wa-breadcrumb-item
      onclick={() =>
        navigate(`/admin/organizers/${contest.ownership.organizerId}/contests`)}
      ><wa-icon name="home"></wa-icon></wa-breadcrumb-item
    >
    <wa-breadcrumb-item onclick={() => navigate(`/admin/contests/${contestId}`)}
      >{contest.name}</wa-breadcrumb-item
    >
    <wa-breadcrumb-item>Results</wa-breadcrumb-item>
  </wa-breadcrumb>

  <h1>Results</h1>

  <ResultsList {contestId} />
{/if}

<style>
  wa-breadcrumb {
    margin-block-end: var(--wa-space-m);
    display: block;
  }
</style>
