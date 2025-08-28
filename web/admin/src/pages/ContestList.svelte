<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import type { Contest } from "@climblive/lib/models";
  import { getContestsByOrganizerQuery } from "@climblive/lib/queries";
  import { format } from "date-fns";
  import { Link, navigate } from "svelte-routing";

  interface Props {
    organizerId: number;
  }

  let { organizerId }: Props = $props();

  const contestsQuery = $derived(getContestsByOrganizerQuery(organizerId));

  let contests = $derived($contestsQuery.data);

  const [drafts, ongoing, upcoming, past] = $derived.by(() => {
    const now = new Date();

    const drafts = contests?.filter(({ timeBegin }) => {
      return !timeBegin;
    });

    const ongoing = contests?.filter(({ timeBegin, timeEnd }) => {
      return timeBegin && timeEnd && now >= timeBegin && now < timeEnd;
    });

    const upcoming = contests?.filter(({ timeBegin }) => {
      return timeBegin && timeBegin > now;
    });

    const past = contests?.filter(({ timeEnd }) => {
      return timeEnd && now > timeEnd;
    });

    return [
      drafts,
      ongoing?.sort(sortContests),
      upcoming?.sort(sortContests),
      past?.sort(sortContests),
    ];
  });

  const sortContests = (c1: Contest, c2: Contest) => {
    if (c1.timeBegin && c2.timeBegin) {
      return c2.timeBegin.getTime() - c1.timeBegin.getTime();
    }

    return 0;
  };

  const columns: ColumnDefinition<Contest>[] = [
    {
      label: "Name",
      mobile: true,
      render: renderName,
      width: "1fr",
    },
    {
      label: "Start time",
      mobile: true,
      render: renderTimeBegin,
      width: "max-content",
    },
    {
      label: "End time",
      mobile: false,
      render: renderTimeEnd,
      width: "max-content",
    },
  ];
</script>

{#snippet renderName({ id, name }: Contest)}
  <Link to="contests/{id}">{name}</Link>
{/snippet}

{#snippet renderTimeBegin({ timeBegin }: Contest)}
  {#if timeBegin}
    {format(timeBegin, "yyyy-MM-dd HH:mm")}
  {/if}
{/snippet}

{#snippet renderTimeEnd({ timeEnd }: Contest)}
  {#if timeEnd}
    {format(timeEnd, "yyyy-MM-dd HH:mm")}
  {/if}
{/snippet}

<h2>Contests</h2>
<wa-button
  variant="brand"
  onclick={() => navigate(`organizers/${organizerId}/contests/new`)}
  >New contest</wa-button
>

{#snippet listing(heading: string, contests: Contest[])}
  <h3>{heading}</h3>
  <Table {columns} data={contests} getId={({ id }) => id}></Table>
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
