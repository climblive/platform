<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import ProblemPoints from "@/components/ProblemPoints.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/progress-bar/progress-bar.js";
  import {
    EmptyState,
    HoldColorIndicator,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import {
    type Contest,
    type Problem,
    type ProblemTemplate,
  } from "@climblive/lib/models";
  import {
    createProblemMutation,
    getAllContestsQuery,
    getProblemsQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";

  interface Props {
    contestId: number;
    open: boolean;
  }

  let { contestId, open = $bindable() }: Props = $props();

  let dialog: WaDialog | undefined = $state();
  let selectedContest: Contest | undefined = $state();
  let isCopying = $state(false);
  let copyProgress = $state(0);
  let copyCompleted = $state(false);

  const contestsQuery = $derived(getAllContestsQuery());
  const selectedProblemsQuery = $derived(
    selectedContest ? getProblemsQuery(selectedContest.id) : undefined,
  );

  const availableContests = $derived.by(() => {
    if (!contestsQuery.data) {
      return [];
    }

    return contestsQuery.data
      .filter(({ id, archived }) => id !== contestId && !archived)
      .sort((a, b) => b.created.getTime() - a.created.getTime());
  });

  const problems = $derived(selectedProblemsQuery?.data);

  const createMutation = $derived(createProblemMutation(contestId));

  const handleClose = () => {
    selectedContest = undefined;
    copyProgress = 0;
    copyCompleted = false;
    open = false;
  };

  const handleSelectContest = (contest: Contest) => {
    selectedContest = contest;
  };

  const handleBack = () => {
    selectedContest = undefined;
  };

  const handleCopy = async () => {
    if (!problems || problems.length === 0) {
      return;
    }

    isCopying = true;
    copyProgress = 0;
    copyCompleted = false;
    let failCount = 0;

    for (let i = 0; i < problems.length; i++) {
      const {
        number,
        holdColorPrimary,
        holdColorSecondary,
        description,
        zone1Enabled,
        zone2Enabled,
        pointsZone1,
        pointsZone2,
        pointsTop,
        flashBonus,
      } = problems[i];

      const template: ProblemTemplate = {
        number,
        holdColorPrimary,
        holdColorSecondary,
        description,
        zone1Enabled,
        zone2Enabled,
        pointsZone1,
        pointsZone2,
        pointsTop,
        flashBonus,
      };

      try {
        await createMutation.mutateAsync(template);
      } catch {
        failCount++;
      }

      copyProgress = ((i + 1) / problems.length) * 100;
    }

    isCopying = false;
    copyCompleted = true;

    if (failCount > 0) {
      toastError(
        `Failed to copy ${failCount} problem${failCount > 1 ? "s" : ""}`,
      );
    }
  };

  const contestColumns: ColumnDefinition<Contest>[] = [
    {
      label: "Contest",
      mobile: true,
      render: renderContestName,
      width: "1fr",
    },
    {
      mobile: true,
      render: renderSelectButton,
      align: "right",
      width: "max-content",
    },
  ];

  const problemColumns: ColumnDefinition<Problem>[] = [
    {
      label: "Number",
      mobile: true,
      render: renderNumberAndColor,
      width: "minmax(max-content, 3fr)",
    },
    {
      label: "Points",
      mobile: true,
      render: renderProblemPoints,
      width: "max-content",
    },
  ];
</script>

{#snippet renderContestName(contest: Contest)}
  {contest.name}
{/snippet}

{#snippet renderContestStartTime(contest: Contest)}
  {#if contest.timeBegin}
    {contest.timeBegin.toLocaleDateString()}
  {:else}
    -
  {/if}
{/snippet}

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

{#snippet renderProblemPoints(
  { pointsZone1, pointsZone2, pointsTop }: Problem,
  mobile: boolean,
)}
  <ProblemPoints {pointsZone1} {pointsZone2} {pointsTop} {mobile} />
{/snippet}

{#snippet renderSelectButton(contest: Contest)}
  <wa-button
    size="small"
    appearance="plain"
    onclick={() => handleSelectContest(contest)}
  >
    Select
    <wa-icon name="arrow-right" slot="end"></wa-icon>
  </wa-button>
{/snippet}

<wa-dialog bind:this={dialog} label="Copy problems" {open}>
  {#if !selectedContest}
    {#if contestsQuery.isPending}
      <Loader />
    {:else if availableContests.length === 0}
      <EmptyState
        title="No contests available"
        description="There are no other contests to copy problems from."
      />
    {:else}
      <Table
        columns={contestColumns}
        data={availableContests}
        getId={({ id }) => id}
        hideHeader={true}
      ></Table>
    {/if}
    <wa-button slot="footer" appearance="plain" onclick={handleClose}>
      Cancel
    </wa-button>
  {:else}
    {#if selectedProblemsQuery?.isPending}
      <Loader />
    {:else if problems && problems.length > 0}
      {#if isCopying || copyCompleted}
        <wa-progress-bar value={copyProgress}></wa-progress-bar>
      {:else}
        <Table
          columns={problemColumns}
          data={problems}
          getId={({ id }) => id}
          hideHeader={true}
        />
      {/if}
    {:else}
      <EmptyState
        title="No problems found"
        description="This contest has no problems to copy."
      />
    {/if}
    {#if copyCompleted}
      <wa-button slot="footer" variant="success" onclick={handleClose}>
        <wa-icon name="check" slot="start"></wa-icon>
        Close
      </wa-button>
    {:else}
      <wa-button
        slot="footer"
        appearance="plain"
        onclick={handleBack}
        disabled={isCopying}
      >
        <wa-icon name="arrow-left" slot="start"></wa-icon>
        Back
      </wa-button>
      <wa-button
        slot="footer"
        variant="accent"
        onclick={handleCopy}
        disabled={!problems || problems.length === 0 || isCopying}
        loading={isCopying}
      >
        Copy {problems?.length || 0} problem{problems?.length !== 1 ? "s" : ""}
        <wa-icon slot="start" name="copy"></wa-icon>
      </wa-button>
    {/if}
  {/if}
</wa-dialog>

<style>
  .number {
    display: flex;
    align-items: center;
    gap: var(--wa-space-s);
  }
</style>
