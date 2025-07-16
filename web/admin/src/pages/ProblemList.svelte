<script lang="ts">
  import {
    HoldColorIndicator,
    Table,
    TableCell,
    TableRow,
  } from "@climblive/lib/components";
  import { getProblemsQuery } from "@climblive/lib/queries";
  import { navigate } from "svelte-routing";
  import DeleteProblem from "./DeleteProblem.svelte";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const problemsQuery = $derived(getProblemsQuery(contestId));

  let problems = $derived($problemsQuery.data);
</script>

<section>
  <Table columns={["Number", "Color", "Points", ""]}>
    {#if problems}
      {#each problems as problem (problem.id)}
        <TableRow>
          <TableCell>№ {problem.number}</TableCell>
          <TableCell
            ><HoldColorIndicator
              --height="1.25rem"
              --width="1.25rem"
              primary={problem.holdColorPrimary}
              secondary={problem.holdColorSecondary}
            /></TableCell
          >
          <TableCell>{problem.pointsTop}</TableCell>
          <TableCell align="right">
            <wa-button
              onclick={() => navigate(`/admin/problems/${problem.id}/edit`)}
              name="pencil"
              label="Edit"
            ></wa-button>
            <DeleteProblem problemId={problem.id}>
              {#snippet children({ deleteProblem })}
                <wa-button
                  onclick={deleteProblem}
                  name="trash"
                  label={`Delete problem ${problem.id}`}
                ></wa-button>
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
    gap: var(--wa-space-xs);
  }
</style>
