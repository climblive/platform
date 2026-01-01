<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import RelativeTime from "@/components/RelativeTime.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import {
    EmptyState,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
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

  let showArchived = $state(false);

  const contestsQuery = $derived(
    getContestsByOrganizerQuery(organizerId ?? 0, {
      enabled: organizerId !== undefined,
    }),
  );
  const allContestsQuery = $derived(
    getAllContestsQuery({ enabled: organizerId === undefined }),
  );

  const contests = $derived(
    organizerId === undefined ? allContestsQuery.data : contestsQuery.data,
  );

  const [ongoing, upcoming, past, archived] = $derived.by(() => {
    const now = new Date();

    const ongoing: Contest[] = [];
    const upcoming: Contest[] = [];
    const past: Contest[] = [];
    const archived: Contest[] = [];

    contests?.forEach((contest) => {
      const { timeBegin, timeEnd } = contest;

      if (contest.archived) {
        archived.push(contest);
      } else if (timeBegin && timeEnd && now >= timeBegin && now < timeEnd) {
        ongoing.push(contest);
      } else if (timeEnd && now > timeEnd) {
        past.push(contest);
      } else {
        upcoming.push(contest);
      }
    });

    ongoing?.sort(sortContests);
    upcoming?.sort(sortContests).reverse();
    past?.sort(sortContests);
    archived?.sort(sortContests);

    return [ongoing, upcoming, past, archived];
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
      label: "Registered",
      mobile: false,
      render: renderRegisteredContenders,
      width: "max-content",
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

  const handleToggleArchive = () => {
    showArchived = !showArchived;
  };

  const numberOfUnarchivedContests = $derived(
    [ongoing, upcoming, past].reduce(
      (partialSum, a) => partialSum + a.length,
      0,
    ),
  );
</script>

{#snippet renderName({ id, name }: Contest)}
  <Link to="contests/{id}">{name}</Link>
{/snippet}

{#snippet renderRegisteredContenders({ registeredContenders }: Contest)}
  {registeredContenders}
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

{#snippet createButton(className?: string)}
  <wa-button
    class={className}
    variant="neutral"
    onclick={() => navigate(`organizers/${organizerId}/contests/new`)}
    >Create new contest</wa-button
  >
{/snippet}

<h2>Contests</h2>
{#if numberOfUnarchivedContests > 0}
  {@render createButton("create-contest-button")}
{/if}

{#snippet listing(
  heading: string,
  contests: Contest[],
  showSummary: boolean = false,
)}
  {@const totalRegistered = contests.reduce(
    (sum, c) => sum + c.registeredContenders,
    0,
  )}
  {@const averageValue =
    contests.length > 0 ? totalRegistered / contests.length : 0}
  {@const averageRegistered = Math.floor(averageValue)}
  <h3>{heading} ({contests.length})</h3>
  {#if showSummary}
    <p class="contest-summary">
      A total of {totalRegistered}
      {totalRegistered === 1 ? "contender has" : "contenders have"} participated in
      {contests.length}
      {contests.length === 1 ? "contest" : "contests"} averaging {averageRegistered}
      {averageRegistered === 1 ? "contender" : "contenders"} per contest.
    </p>
  {/if}
  <Table {columns} data={contests} getId={({ id }) => id}></Table>
{/snippet}

{#if !ongoing || !upcoming || !past || !archived}
  <Loader />
{:else}
  {#if ongoing?.length}
    {@render listing("Ongoing", ongoing)}
  {/if}

  {#if upcoming?.length}
    {@render listing("Upcoming", upcoming)}
  {/if}

  {#if past?.length}
    {@render listing("Past", past, true)}
  {/if}

  {#if contests && numberOfUnarchivedContests === 0}
    <EmptyState
      title="No contests yet"
      description="Create a contest to get started with your first event."
    >
      {#snippet actions()}
        {@render createButton()}
      {/snippet}
    </EmptyState>
  {/if}

  {#if showArchived && archived.length > 0}
    {@render listing("Archived", archived)}
  {/if}

  {#if archived?.length}
    <wa-button
      size="small"
      appearance="plain"
      variant="brand"
      onclick={handleToggleArchive}
      class="toggle-archived-button"
    >
      {#if showArchived}
        Hide archived contests
      {:else}
        Show archived contests ({archived.length})
      {/if}
    </wa-button>
  {/if}
{/if}

<style>
  .create-contest-button {
    margin-block-end: var(--wa-space-m);
  }

  .toggle-archived-button {
    display: block;
    margin-block-start: var(--wa-space-m);
  }

  .contest-summary {
    color: var(--wa-color-text-quiet);
    font-size: var(--wa-font-size-s);
    margin-block-start: var(--wa-space-xs);
    margin-block-end: var(--wa-space-m);
  }
</style>
