<script lang="ts">
  import type { Tick } from "@climblive/lib/models";
  import {
    getProblemsQuery,
    getTicksByContenderQuery,
  } from "@climblive/lib/queries";
  import Problem from "./Problem.svelte";

  type Props = {
    contestId: number;
    contenderId: number;
  };

  const { contestId, contenderId }: Props = $props();

  const problemsQuery = $derived(getProblemsQuery(contestId));
  const ticksQuery = $derived(getTicksByContenderQuery(contenderId));

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

      <Problem {problem} {tick} {contenderId} />
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
