<script lang="ts">
  import { Table } from "@climblive/lib/components";
  import { getCompClassesQuery } from "@climblive/lib/queries";
  import { format } from "date-fns";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const compClassesQuery = getCompClassesQuery(contestId);

  let compClasses = $derived($compClassesQuery.data);
</script>

<section>
  <Table columns={["Name", "Start time", "End time"]}>
    {#if compClasses}
      {#each compClasses as compClass (compClass.id)}
        <span>{compClass.name}</span>
        <span>
          {format(compClass.timeBegin, "yyyy-MM-dd HH:mm")}
        </span>
        <span>
          {format(compClass.timeEnd, "yyyy-MM-dd HH:mm")}
        </span>
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
