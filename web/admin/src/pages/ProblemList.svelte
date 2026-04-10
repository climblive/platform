<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import { type WaSelectEvent } from "@awesome.me/webawesome";
  import WaDropdownItem from "@awesome.me/webawesome/dist/components/dropdown-item/dropdown-item.js";
  import "@awesome.me/webawesome/dist/components/dropdown/dropdown.js";
  import {
    EmptyState,
    HoldColorIndicator,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import { type Problem, type ProblemID } from "@climblive/lib/models";
  import {
    getContestsByOrganizerQuery,
    getProblemsQuery,
    getTicksByContestQuery,
  } from "@climblive/lib/queries";
  import { isDefined } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";
  import CopyProblems from "./CopyProblems.svelte";
  import DeleteProblem from "./DeleteProblem.svelte";

  const maxProblems = 100;

  interface Props {
    contestId: number;
    organizerId: number;
    tableLimit: number | undefined;
  }

  let { contestId, organizerId, ...props }: Props = $props();
  let tableLimit = $state<number | undefined>(props.tableLimit);

  let copyProblemsOpen = $state(false);

  const problemsQuery = $derived(getProblemsQuery(contestId));
  const ticksQuery = $derived(getTicksByContestQuery(contestId));
  const contestsQuery = $derived(getContestsByOrganizerQuery(organizerId));

  const ascentsByProblem = $derived.by(() => {
    const ascentsByProblem = new Map<ProblemID, number>();

    if (ticksQuery.data === undefined) {
      return ascentsByProblem;
    }

    for (const { problemId } of ticksQuery.data) {
      const ascents = ascentsByProblem.get(problemId);

      if (ascents !== undefined) {
        ascentsByProblem.set(problemId, ascents + 1);
      } else {
        ascentsByProblem.set(problemId, 1);
      }
    }

    return ascentsByProblem;
  });
  const contests = $derived(contestsQuery.data);

  const limitReached = $derived(
    problemsQuery.data !== undefined &&
      problemsQuery.data.length >= maxProblems,
  );

  type ProblemWithAscents = Problem & { ascents: number };

  const sortedProblemsWithAscents = $derived.by(() => {
    if (problemsQuery.data === undefined) {
      return undefined;
    }

    const problems = problemsQuery.data.map<ProblemWithAscents>((problem) => ({
      ...problem,
      ascents: ascentsByProblem.get(problem.id) ?? 0,
    }));

    return problems?.sort((p1, p2) => p1.number - p2.number);
  });

  const showAll = () => {
    tableLimit = undefined;
  };

  const columns: ColumnDefinition<ProblemWithAscents>[] = [
    {
      label: "Number",
      mobile: true,
      render: renderNumberAndColor,
      width: "minmax(max-content, 3fr)",
    },
    {
      label: "Zones",
      mobile: false,
      render: renderZones,
      width: "max-content",
    },
    {
      label: "Points",
      mobile: true,
      render: renderPoints,
      width: "max-content",
    },
    {
      label: "Flash",
      mobile: false,
      render: renderFlashBonus,
      width: "max-content",
    },
    {
      label: "Tops",
      mobile: false,
      render: renderAscents,
      align: "right",
      width: "max-content",
    },
    {
      mobile: true,
      render: renderControls,
      align: "right",
      width: "max-content",
    },
  ];
</script>

{#snippet renderNumberAndColor({
  number,
  holdColorPrimary,
  holdColorSecondary,
}: Problem)}
  <div class="number">
    <HoldColorIndicator
      --height="1.25rem"
      --width="1.25rem"
      primary={holdColorPrimary}
      secondary={holdColorSecondary}
    />
    № {number}
  </div>
{/snippet}

{#snippet renderPoints(
  { pointsZone1, pointsZone2, pointsTop, flashBonus }: ProblemWithAscents,
  mobile: boolean,
)}
  {@const values = [pointsZone1, pointsZone2, pointsTop].filter(isDefined)}

  {#if mobile}
    {@const min = Math.min(...values)}
    {@const max = Math.max(...values)}

    {[min, max + (flashBonus ?? 0)].join(" - ")} pts
  {:else}
    {values.join(" / ")} pts
  {/if}
{/snippet}

{#snippet renderFlashBonus({ flashBonus }: ProblemWithAscents)}
  {#if flashBonus}
    +{flashBonus} pts
  {:else}
    -
  {/if}
{/snippet}

{#snippet renderZones({ zone1Enabled, zone2Enabled }: ProblemWithAscents)}
  {#if zone1Enabled && zone2Enabled}
    Z1 + Z2
  {:else if zone1Enabled}
    Z1
  {:else}
    -
  {/if}
{/snippet}

{#snippet renderControls({ id }: ProblemWithAscents)}
  <DeleteProblem problemId={id}>
    {#snippet children({ deleteProblem })}
      <wa-dropdown
        onwa-select={(event: WaSelectEvent) => {
          if ((event.detail.item as WaDropdownItem).value === "delete") {
            deleteProblem();
          } else {
            navigate(`/admin/problems/${id}/edit`);
          }
        }}
      >
        <wa-button slot="trigger" size="small" appearance="plain">
          <wa-icon name="ellipsis-vertical" label="Actions"></wa-icon>
        </wa-button>
        <wa-dropdown-item value="edit">
          <wa-icon slot="icon" name="pencil"></wa-icon>
          Edit
        </wa-dropdown-item>
        <wa-dropdown-item value="delete" variant="danger">
          <wa-icon slot="icon" name="trash"></wa-icon>
          Delete
        </wa-dropdown-item>
      </wa-dropdown>
    {/snippet}
  </DeleteProblem>
{/snippet}

{#snippet renderAscents({ ascents }: ProblemWithAscents)}
  {ascents}
{/snippet}

{#snippet createButton()}
  <wa-button
    variant="neutral"
    appearance="accent"
    onclick={() => navigate(`contests/${contestId}/new-problem`)}
    disabled={limitReached}>Create problem</wa-button
  >
  {#if limitReached}
    <wa-callout variant="warning">
      You have reached the maximum of {maxProblems} problems per contest.
    </wa-callout>
  {/if}
{/snippet}

<p class="copy">
  Problems refer to the boulder problems that the contenders will attempt during
  the contest, each of which can have its own point value.
</p>

<section>
  {#if sortedProblemsWithAscents === undefined}
    <Loader />
  {:else if sortedProblemsWithAscents.length > 0}
    {@render createButton()}
    <Table
      {columns}
      data={tableLimit
        ? sortedProblemsWithAscents.slice(0, tableLimit)
        : sortedProblemsWithAscents}
      getId={({ id }) => id}
    ></Table>
  {:else}
    <EmptyState
      title="No problems yet"
      description="Create boulder problems that contenders will attempt during the contest."
    >
      {#snippet actions()}
        <div class="actions">
          {@render createButton()}

          {#if contests && contests.length > 1}
            <wa-button
              onclick={() => {
                copyProblemsOpen = true;
              }}
              appearance="outlined"
              variant="neutral"
            >
              Copy from another contest
              <wa-icon name="copy" slot="start"></wa-icon>
            </wa-button>
          {/if}
        </div>
      {/snippet}
    </EmptyState>
  {/if}

  {#if sortedProblemsWithAscents !== undefined && tableLimit !== undefined && tableLimit < sortedProblemsWithAscents.length}
    <wa-button class="show-more" appearance="plain" onclick={showAll}
      >Show all</wa-button
    >
  {/if}
</section>

<CopyProblems {organizerId} {contestId} bind:open={copyProblemsOpen} />

<style>
  .number {
    display: flex;
    align-items: center;
    gap: var(--wa-space-s);
  }

  section {
    display: flex;
    flex-direction: column;
    align-items: start;
    gap: var(--wa-space-m);
  }

  wa-button.show-more {
    align-self: center;
  }

  .copy {
    color: var(--wa-color-text-quiet);
  }

  .actions {
    display: flex;
    gap: var(--wa-space-xs);
    flex-wrap: wrap;
    justify-content: center;
  }
</style>
