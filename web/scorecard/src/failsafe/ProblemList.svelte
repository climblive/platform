<script lang="ts">
  import type { Tick } from "@climblive/lib/models";
  import {
    getProblemsQuery,
    getTicksByCompetitorQuery,
  } from "@climblive/lib/queries";
  import Problem from "./Problem.svelte";

  type Props = {
    competitionId: string;
    competitorId: number;
    enableAttempts: boolean;
  };

  const { competitionId, competitorId, enableAttempts }: Props = $props();

  const problemsQuery = $derived(getProblemsQuery(competitionId));
  const ticksQuery = $derived(getTicksByCompetitorQuery(competitorId));

  let problems = $derived(
    problemsQuery.data
      ? [...problemsQuery.data]?.sort((a, b) => a.number - b.number)
      : undefined,
  );
  let ticksByProblem = $derived<Map<number, Tick>>(
    new Map(ticksQuery.data?.map((tick) => [tick.problemId, tick]) ?? []),
  );
</script>

{#if problems}
  <section>
    {#each problems as problem (problem.id)}
      {@const tick = ticksByProblem.get(problem.id)}

      <Problem {problem} {tick} {competitorId} {enableAttempts} />
    {/each}
  </section>
{/if}

<style>
  section {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
    margin-top: var(--wa-space-m);
  }
</style>
