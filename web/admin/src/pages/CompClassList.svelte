<script lang="ts">
  import { Table, TableCell, TableRow } from "@climblive/lib/components";
  import { getCompClassesQuery } from "@climblive/lib/queries";
  import { format } from "date-fns";
  import DeleteCompClass from "./DeleteCompClass.svelte";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const compClassesQuery = getCompClassesQuery(contestId);

  let compClasses = $derived($compClassesQuery.data);
</script>

<section>
  <Table columns={["Name", "Start time", "End time", ""]}>
    {#if compClasses}
      {#each compClasses as compClass (compClass.id)}
        <TableRow>
          <TableCell>{compClass.name}</TableCell>
          <TableCell>
            {format(compClass.timeBegin, "yyyy-MM-dd HH:mm")}
          </TableCell>
          <TableCell>
            {format(compClass.timeEnd, "yyyy-MM-dd HH:mm")}
          </TableCell>
          <TableCell align="right">
            <DeleteCompClass compClassId={compClass.id}>
              {#snippet children({ deleteCompClass })}
                <sl-icon-button
                  onclick={deleteCompClass}
                  name="trash"
                  label={`Delete comp class ${compClass.id}`}
                ></sl-icon-button>
              {/snippet}
            </DeleteCompClass>
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
