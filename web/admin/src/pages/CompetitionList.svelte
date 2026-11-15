<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import RelativeTime from "@/components/RelativeTime.svelte";
  import { type WaSelectEvent } from "@awesome.me/webawesome";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/divider/divider.js";
  import WaDropdownItem from "@awesome.me/webawesome/dist/components/dropdown-item/dropdown-item.js";
  import "@awesome.me/webawesome/dist/components/dropdown/dropdown.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/tooltip/tooltip.js";
  import {
    EmptyState,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import type { Competition } from "@climblive/lib/models";
  import {
    getAllCompetitionsQuery,
    getCompetitionsByOrganizerQuery,
    getSelfQuery,
  } from "@climblive/lib/queries";
  import { getCountryName, getFlag, SyncedTime } from "@climblive/lib/utils";
  import { format, sub } from "date-fns";
  import { onMount } from "svelte";
  import { Link, navigate } from "svelte-routing";

  import ArchiveCompetition from "./ArchiveCompetition.svelte";
  import DuplicateCompetition from "./DuplicateCompetition.svelte";
  import RestoreCompetition from "./RestoreCompetition.svelte";
  import TransferCompetition from "./TransferCompetition.svelte";

  const maxCompetitionsPerWeek = 10;
  const defaultTableLimit = 10;
  const createButtonId = $props.id();

  interface Props {
    organizerId: number | undefined;
  }

  let { organizerId }: Props = $props();

  let showArchived = $state(false);
  let expandedListings = $state<string[]>([]);
  let selectedAction = $state<
    | { type: "duplicate" | "transfer" | "archive"; competition: Competition }
    | undefined
  >();

  const time = new SyncedTime(60_000);

  onMount(() => {
    time.start();

    return () => time.stop();
  });

  const selfQuery = $derived(getSelfQuery());

  const competitionsQuery = $derived(
    getCompetitionsByOrganizerQuery(organizerId ?? 0, {
      enabled: organizerId !== undefined,
    }),
  );
  const allCompetitionsQuery = $derived(
    getAllCompetitionsQuery({ enabled: organizerId === undefined }),
  );

  const self = $derived(selfQuery.data);
  const competitions = $derived(
    organizerId === undefined
      ? allCompetitionsQuery.data
      : competitionsQuery.data,
  );

  const [ongoing, upcoming, past, archived] = $derived.by(() => {
    const now = time.current;

    const ongoing: Competition[] = [];
    const upcoming: Competition[] = [];
    const past: Competition[] = [];
    const archived: Competition[] = [];

    competitions?.forEach((competition) => {
      const { timeBegin, timeEnd } = competition;

      if (competition.archivedAt != null) {
        archived.push(competition);
      } else if (timeBegin && timeEnd && now >= timeBegin && now < timeEnd) {
        ongoing.push(competition);
      } else if (timeEnd && now > timeEnd) {
        past.push(competition);
      } else {
        upcoming.push(competition);
      }
    });

    ongoing?.sort(sortCompetitions);
    upcoming?.sort(sortCompetitions).reverse();
    past?.sort(sortCompetitions);
    archived?.sort(sortCompetitions);

    return [ongoing, upcoming, past, archived];
  });

  const sortCompetitions = (c1: Competition, c2: Competition) => {
    if (c1.timeBegin && c2.timeBegin) {
      return c2.timeBegin.getTime() - c1.timeBegin.getTime();
    }

    return 0;
  };

  const columns: ColumnDefinition<Competition>[] = [
    {
      label: "",
      mobile: true,
      render: renderFlag,
      width: "max-content",
    },
    {
      label: "Name",
      mobile: true,
      render: renderName,
      width: "2fr",
    },
    {
      label: "Location",
      mobile: false,
      render: renderLocation,
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
    {
      label: "Registered",
      mobile: false,
      render: renderRegisteredCompetitors,
      width: "max-content",
      align: "right",
    },
    {
      mobile: true,
      render: renderControls,
      align: "right",
      width: "max-content",
    },
  ];

  const handleToggleArchive = () => {
    showArchived = !showArchived;
  };

  const expandListing = (heading: string) => {
    expandedListings = [...expandedListings, heading];
  };

  const handleCompetitionAction = (
    event: WaSelectEvent,
    competition: Competition,
  ) => {
    const action = (event.detail.item as WaDropdownItem).value;

    switch (action) {
      case "edit":
        navigate(`/admin/competitions/${competition.id}/edit`);
        break;
      case "results":
        navigate(`/admin/competitions/${competition.id}/results`);
        break;
      case "duplicate":
      case "transfer":
      case "archive":
        selectedAction = { type: action, competition };
        break;
    }
  };

  const numberOfUnarchivedCompetitions = $derived(
    [ongoing, upcoming, past].reduce(
      (partialSum, a) => partialSum + a.length,
      0,
    ),
  );

  const limitReached = $derived.by(() => {
    if (!competitions) {
      return false;
    }

    if (self?.admin) {
      return false;
    }

    const oneWeekAgo = sub(new Date(), { weeks: 1 });
    const recentCount = competitions.filter(
      (c) => c.archivedAt == null && c.created > oneWeekAgo,
    ).length;

    return recentCount >= maxCompetitionsPerWeek;
  });
</script>

{#snippet renderFlag({ id, country }: Competition)}
  <wa-tooltip for="flag-{id}">{getCountryName(country)}</wa-tooltip>
  <span id="flag-{id}" class="flag">{getFlag(country)}</span>
{/snippet}

{#snippet renderName({ id, name }: Competition)}
  <Link to="competitions/{id}">{name}</Link>
{/snippet}

{#snippet renderLocation({ location }: Competition)}
  {location || "-"}
{/snippet}

{#snippet renderRegisteredCompetitors({ registeredCompetitors }: Competition)}
  {registeredCompetitors}
{/snippet}

{#snippet renderTimeBegin({ timeBegin, timeEnd }: Competition, mobile: boolean)}
  {#if timeBegin}
    {#if timeEnd && new Date() > timeEnd}
      {format(timeBegin, "yyyy-MM-dd HH:mm")}
    {:else}
      <RelativeTime format={mobile ? "narrow" : "long"} time={timeBegin} />
    {/if}
  {:else}
    -
  {/if}
{/snippet}

{#snippet renderTimeEnd({ timeEnd }: Competition)}
  {#if timeEnd}
    {format(timeEnd, "yyyy-MM-dd HH:mm")}
  {:else}
    -
  {/if}
{/snippet}

{#snippet renderControls(competition: Competition)}
  {#if competition.archivedAt !== undefined}
    <RestoreCompetition competitionId={competition.id}>
      {#snippet children({ restoreCompetition })}
        <wa-dropdown
          onwa-select={(event: WaSelectEvent) => {
            if ((event.detail.item as WaDropdownItem).value === "restore") {
              restoreCompetition();
            }
          }}
        >
          <wa-button slot="trigger" size="s" appearance="plain">
            <wa-icon name="ellipsis-vertical" label="Actions"></wa-icon>
          </wa-button>
          <wa-dropdown-item value="restore">
            <wa-icon slot="icon" name="rotate-left"></wa-icon>
            Restore
          </wa-dropdown-item>
        </wa-dropdown>
      {/snippet}
    </RestoreCompetition>
  {:else}
    <wa-dropdown
      onwa-select={(event: WaSelectEvent) =>
        handleCompetitionAction(event, competition)}
    >
      <wa-button slot="trigger" size="s" appearance="plain">
        <wa-icon name="ellipsis-vertical" label="Actions"></wa-icon>
      </wa-button>
      <wa-dropdown-item value="edit">
        <wa-icon slot="icon" name="pencil"></wa-icon>
        Edit
      </wa-dropdown-item>
      <wa-dropdown-item value="results">
        <wa-icon slot="icon" name="ranking-star"></wa-icon>
        View results
      </wa-dropdown-item>

      <wa-divider></wa-divider>

      <wa-dropdown-item value="duplicate">
        <wa-icon slot="icon" name="copy"></wa-icon>
        Duplicate
      </wa-dropdown-item>
      <wa-dropdown-item
        value="transfer"
        disabled={!self?.organizers.some(
          ({ id }) => id !== competition.ownership.organizerId,
        )}
      >
        <wa-icon slot="icon" name="arrow-right"></wa-icon>
        Transfer
      </wa-dropdown-item>
      <wa-dropdown-item value="archive" variant="danger">
        <wa-icon slot="icon" name="box-archive"></wa-icon>
        Archive
      </wa-dropdown-item>
    </wa-dropdown>
  {/if}
{/snippet}

{#snippet createButton(className?: string)}
  {#if organizerId}
    <wa-button
      id={createButtonId}
      class={className}
      variant="neutral"
      disabled={limitReached}
      onclick={() => navigate(`organizers/${organizerId}/competitions/new`)}
      >Create new competition</wa-button
    >
    {#if limitReached}
      <wa-tooltip for={createButtonId} placement="top-start"
        >Limit of {maxCompetitionsPerWeek} competitions per week reached.</wa-tooltip
      >
    {/if}
  {/if}
{/snippet}

<h2>Competitions</h2>
{#if numberOfUnarchivedCompetitions > 0}
  {@render createButton("create-competition-button")}
{/if}

{#snippet listing(
  heading: string,
  competitions: Competition[],
  showSummary: boolean = false,
)}
  {@const totalRegistered = competitions.reduce(
    (sum, c) => sum + c.registeredCompetitors,
    0,
  )}
  {@const averageRegistered = Math.floor(
    competitions.length > 0 ? totalRegistered / competitions.length : 0,
  )}

  <h3>{heading} ({competitions.length})</h3>
  {#if showSummary}
    <p class="competition-summary">
      A total of {totalRegistered}
      {totalRegistered === 1 ? "competitor has" : "competitors have"} participated
      in {competitions.length}
      {competitions.length === 1 ? "competition" : "competitions"} averaging {averageRegistered}
      {averageRegistered === 1 ? "competitor" : "competitors"} per competition.
    </p>
  {/if}
  <Table
    {columns}
    data={expandedListings.includes(heading)
      ? competitions
      : competitions.slice(0, defaultTableLimit)}
    getId={({ id }) => id}
  ></Table>
  {#if !expandedListings.includes(heading) && competitions.length > defaultTableLimit}
    <wa-button
      class="show-more"
      appearance="plain"
      onclick={() => expandListing(heading)}>Show all</wa-button
    >
  {/if}
{/snippet}

{#if competitionsQuery.isLoading || allCompetitionsQuery.isLoading || !ongoing || !upcoming || !past || !archived}
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

  {#if competitions && numberOfUnarchivedCompetitions === 0}
    <EmptyState
      title="No competitions yet"
      description="Create a competition to get started with your first event."
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
      size="s"
      appearance="plain"
      variant="brand"
      onclick={handleToggleArchive}
      class="toggle-archived-button"
    >
      {#if showArchived}
        Hide archived competitions
      {:else}
        Show archived competitions ({archived.length})
      {/if}
    </wa-button>
  {/if}
{/if}

{#if selectedAction?.type === "duplicate"}
  <DuplicateCompetition
    competitionId={selectedAction.competition.id}
    competitionName={selectedAction.competition.name}
    open
    onClose={() => (selectedAction = undefined)}
  />
{:else if selectedAction?.type === "transfer"}
  <TransferCompetition
    competitionId={selectedAction.competition.id}
    competitionName={selectedAction.competition.name}
    organizerId={selectedAction.competition.ownership.organizerId}
    open
    onClose={() => (selectedAction = undefined)}
  />
{:else if selectedAction?.type === "archive"}
  <ArchiveCompetition
    competitionId={selectedAction.competition.id}
    competitionName={selectedAction.competition.name}
    organizerId={selectedAction.competition.ownership.organizerId}
    open
    onClose={() => (selectedAction = undefined)}
  />
{/if}

<style>
  .create-competition-button {
    margin-block-end: var(--wa-space-m);
  }

  .toggle-archived-button {
    display: block;
    margin-block-start: var(--wa-space-m);
  }

  wa-button.show-more {
    display: block;
    margin: var(--wa-space-m) auto 0;
  }

  .competition-summary {
    color: var(--wa-color-text-quiet);
    font-size: var(--wa-font-size-s);
    margin-block-start: var(--wa-space-xs) var(--wa-space-m);
  }

  .flag {
    font-size: var(--wa-font-size-larger);
  }
</style>
