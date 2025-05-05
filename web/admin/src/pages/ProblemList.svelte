<script lang="ts">
  import { Table, TableCell, TableRow } from "@climblive/lib/components";
  import { getProblemsQuery } from "@climblive/lib/queries";
  import "@shoelace-style/shoelace/dist/components/icon-button/icon-button.js";
  import { navigate } from "svelte-routing";
  import DeleteProblem from "./DeleteProblem.svelte";

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
          <TableCell>№ {problem.number}</TableCell>
          <TableCell>{problem.holdColorPrimary}</TableCell>
          <TableCell>{problem.pointsTop}</TableCell>
          <TableCell align="right">
            <sl-icon-button
              onclick={() => navigate(`/admin/problems/${problem.id}/edit`)}
              name="pencil"
              label="Edit"
            ></sl-icon-button>
            <DeleteProblem problemID={problem.id}>
              {#snippet children({ onDelete })}
                <sl-icon-button onclick={onDelete} name="trash" label="Remove"
                ></sl-icon-button>
              {/snippet}
            </DeleteProblem>
          </TableCell>
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
