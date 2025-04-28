<script lang="ts">
  import { Table } from "@climblive/lib/components";
  import { getProblemsQuery } from "@climblive/lib/queries";
  import { navigate } from "svelte-routing";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const problemsQuery = getProblemsQuery(contestId);

  let problems = $derived($problemsQuery.data);
</script>

<section>
  <Table columns={["Number", "Color", "Points"]}>
    {#if problems}
      {#each problems as problem (problem.id)}
        <span>№ {problem.number}</span>
        <span>{problem.holdColorPrimary}</span>
        <sl-button
          onclick={() => navigate(`/admin/problems/${problem.id}/edit`)}
          >Edit</sl-button
        >
      {/each}
    {/if}
  </Table>
</section>

<style>
  section {
    display: flex;
    gap: var(--sl-spacing-x-small);
  }
</style>
