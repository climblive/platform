<script lang="ts">
  import { Table, TableRow } from "@climblive/lib/components";
  import { getProblemsQuery } from "@climblive/lib/queries";
  import "@shoelace-style/shoelace/dist/components/icon-button/icon-button.js";
  import { navigate } from "svelte-routing";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const problemsQuery = getProblemsQuery(contestId);

  let problems = $derived($problemsQuery.data);
</script>

<section>
  <Table columns={["Number", "Color", "Points", ""]}>
    {#if problems}
      {#each problems as problem (problem.id)}
        <TableRow>
          <span>№ {problem.number}</span>
          <span>{problem.holdColorPrimary}</span>
          <span>{problem.pointsTop}</span>
          <sl-icon-button
            onclick={() => navigate(`/admin/problems/${problem.id}/edit`)}
            name="pencil"
            label="Edit"
          ></sl-icon-button>
        </TableRow>
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
