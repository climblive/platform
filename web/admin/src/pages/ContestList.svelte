<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import RelativeTime from "@/components/RelativeTime.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/tag/tag.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import type { Contest } from "@climblive/lib/models";
  import {
    getAllContestsQuery,
    getContestsByOrganizerQuery,
  } from "@climblive/lib/queries";
  import { format } from "date-fns";
  import { Link, navigate } from "svelte-routing";

  interface Props {
    organizerId: number | undefined;
  }

  let { organizerId }: Props = $props();

  const contestsQuery = $derived(
    getContestsByOrganizerQuery(organizerId ?? 0, {
      enabled: organizerId !== undefined,
    }),
  );
  const allContestsQuery = $derived(
    getAllContestsQuery({ enabled: organizerId === undefined }),
  );

  const contests = $derived(
    organizerId === undefined ? $allContestsQuery.data : $contestsQuery.data,
  );

  const [ongoing, upcoming, past] = $derived.by(() => {
    const now = new Date();

    const ongoing: Contest[] = [];
    const upcoming: Contest[] = [];
    const past: Contest[] = [];

    contests?.forEach((contest) => {
      const { timeBegin, timeEnd } = contest;

      if (timeBegin && timeEnd && now >= timeBegin && now < timeEnd) {
        ongoing.push(contest);
      } else if (timeEnd && now > timeEnd) {
        past.push(contest);
      } else {
        upcoming.push(contest);
      }
    });

    ongoing?.sort(sortContests);
    upcoming?.sort(sortContests);
    past?.sort(sortContests);

    return [ongoing, upcoming, past];
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

{#snippet renderTimeBegin({ timeBegin, timeEnd }: Contest)}
  {#if timeBegin}
    {#if timeEnd && new Date() > timeEnd}
      {format(timeBegin, "yyyy-MM-dd HH:mm")}
    {:else}
      <RelativeTime time={timeBegin} />
    {/if}
  {:else}
    -
  {/if}
{/snippet}

{#snippet renderTimeEnd({ timeEnd }: Contest)}
  {#if timeEnd}
    {format(timeEnd, "yyyy-MM-dd HH:mm")}
  {:else}
    -
  {/if}
{/snippet}

<h2>Contests</h2>
<wa-button
  variant="neutral"
  onclick={() => navigate(`organizers/${organizerId}/contests/new`)}
  >Create new contest</wa-button
>

{#snippet listing(heading: string, contests: Contest[])}
  <h3>{heading}</h3>
  <Table {columns} data={contests} getId={({ id }) => id}></Table>
{/snippet}

{#if !ongoing || !upcoming || !past}
  <Loader />
{:else}
  {#if ongoing?.length}
    {@render listing("Ongoing", ongoing)}
  {/if}

  {#if upcoming?.length}
    {@render listing("Upcoming", upcoming)}
  {/if}

  {#if past?.length}
    {@render listing("Past", past)}
  {/if}
{/if}

<style>
  wa-button {
    margin-block-end: var(--wa-space-m);
  }
</style>
