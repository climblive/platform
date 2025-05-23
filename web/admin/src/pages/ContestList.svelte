<script lang="ts">
  import { Table, TableCell, TableRow } from "@climblive/lib/components";
  import type { Contest } from "@climblive/lib/models";
  import { getContestsByOrganizerQuery } from "@climblive/lib/queries";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import { format } from "date-fns";
  import { Link, navigate } from "svelte-routing";

  interface Props {
    organizerId: number;
  }

  let { organizerId }: Props = $props();

  const contestsQuery = getContestsByOrganizerQuery(organizerId);

  let contests = $derived($contestsQuery.data);

  const [drafts, ongoing, upcoming, past] = $derived.by(() => {
    const now = new Date();

    const drafts = contests?.filter(({ timeBegin }) => {
      return !timeBegin;
    });

    const ongoing = contests?.filter(({ timeBegin, timeEnd }) => {
      return timeBegin && timeEnd && timeBegin >= now && timeEnd < now;
    });

    const upcoming = contests?.filter(({ timeBegin }) => {
      return timeBegin && timeBegin > now;
    });

    const past = contests?.filter(({ timeEnd }) => {
      return timeEnd && now > timeEnd;
    });

    return [drafts, ongoing, upcoming, past];
  });
</script>

<sl-button
  size="large"
  variant="primary"
  onclick={() => navigate(`organizers/${organizerId}/contests/new`)}
  >Create</sl-button
>

{#snippet listing(heading: string, contests: Contest[])}
  <h2>{heading}</h2>
  <Table columns={["Name", "Start Time", "End Time"]}>
    {#each contests as contest (contest.id)}
      <TableRow>
        <TableCell>
          <Link to="contests/{contest.id}">{contest.name}</Link>
        </TableCell>
        <TableCell>
          {#if contest.timeBegin}
            {format(contest.timeBegin, "yyyy-MM-dd HH:mm")}
          {/if}
        </TableCell>
        <TableCell>
          {#if contest.timeEnd}
            {format(contest.timeEnd, "yyyy-MM-dd HH:mm")}
          {/if}
        </TableCell>
      </TableRow>
    {/each}
  </Table>
{/snippet}

{#if drafts?.length}
  {@render listing("Drafts", drafts)}
{/if}

{#if ongoing?.length}
  {@render listing("Ongoing", ongoing)}
{/if}

{#if upcoming?.length}
  {@render listing("Upcoming", upcoming)}
{/if}

{#if past?.length}
  {@render listing("Past", past)}
{/if}
