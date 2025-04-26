<script lang="ts">
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
  {#if problems}
    <ul>
      {#each problems as problem (problem.id)}
        <li>
          {problem.number}
          <sl-button
            onclick={() => navigate(`/admin/problems/${problem.id}/edit`)}
            >Edit</sl-button
          >
        </li>
      {/each}
    </ul>
  {/if}
</section>

<style>
  section {
    display: flex;
    gap: var(--sl-spacing-x-small);
  }
</style>
