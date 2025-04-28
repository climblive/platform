<script lang="ts">
  import {
    getContestQuery,
    getScoreEnginesQuery,
    startScoreEngineMutation,
    stopScoreEngineMutation,
  } from "@climblive/lib/queries";
  import "@shoelace-style/shoelace/dist/components/tab-group/tab-group.js";
  import "@shoelace-style/shoelace/dist/components/tab-panel/tab-panel.js";
  import "@shoelace-style/shoelace/dist/components/tab/tab.js";
  import { add } from "date-fns";
  import { Link } from "svelte-routing";
  import CompClassList from "./CompClassList.svelte";
  import ContenderList from "./ContenderList.svelte";
  import CreateCompClass from "./CreateCompClass.svelte";
  import CreateProblem from "./CreateProblem.svelte";
  import ProblemList from "./ProblemList.svelte";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const contestQuery = getContestQuery(contestId);
  const scoreEnginesQuery = getScoreEnginesQuery(contestId);
  const startScoreEngine = startScoreEngineMutation(contestId);
  const stopScoreEngine = stopScoreEngineMutation();

  let contest = $derived($contestQuery.data);
  let scoreEngines = $derived($scoreEnginesQuery.data);
</script>

<main>
  <Link to="/admin/contests">Back to contests</Link>

  {#if contest && scoreEngines}
    <h1>{contest.name}</h1>

    <sl-tab-group>
      <sl-tab slot="nav" panel="contest">Contest</sl-tab>
      <sl-tab slot="nav" panel="contenders">Contenders</sl-tab>
      <sl-tab slot="nav" panel="problems">Problems</sl-tab>

      <sl-tab-panel name="contest">
        <h2>Score Engines</h2>
        {#each scoreEngines as engineInstanceId}
          <div>
            <h3>{engineInstanceId}</h3>
            <sl-button
              variant="danger"
              onclick={() => $stopScoreEngine.mutate(engineInstanceId)}
              loading={$stopScoreEngine.isPending}>Stop</sl-button
            >
          </div>
        {/each}
        <sl-button
          onclick={() =>
            $startScoreEngine.mutate({
              terminatedBy: add(new Date(), { hours: 6 }),
            })}
          loading={$startScoreEngine.isPending}
          disabled={scoreEngines.length > 0}>Start engine</sl-button
        >

        <h2>Classes</h2>
        <CreateCompClass {contestId} />
        <CompClassList {contestId} />
      </sl-tab-panel>

      <sl-tab-panel name="problems">
        <h2>Problems</h2>
        <CreateProblem {contestId} />
        <ProblemList {contestId} />
      </sl-tab-panel>

      <sl-tab-panel name="contenders">
        <h2>Contenders</h2>
        <ContenderList {contestId} />
      </sl-tab-panel>
    </sl-tab-group>
  {/if}
</main>

<style>
  main {
    padding: var(--sl-spacing-medium);
  }
</style>
