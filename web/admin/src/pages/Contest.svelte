<script lang="ts">
  import {
    getContestQuery,
    getScoreEnginesQuery,
    startScoreEngineMutation,
    stopScoreEngineMutation,
  } from "@climblive/lib/queries";
  import { add } from "date-fns";
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
  {#if contest && scoreEngines}
    <h1>{contest.name}</h1>

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

    <h2>Problems</h2>
    <CreateProblem {contestId} />
    <ProblemList {contestId} />

    <h2>Contenders</h2>
    <ContenderList {contestId} />
  {/if}
</main>

<style>
  main {
    padding: var(--sl-spacing-medium);
  }
</style>
