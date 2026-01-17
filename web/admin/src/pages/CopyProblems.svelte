<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/progress-bar/progress-bar.js";
  import {
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
    getAllContestsQuery,
    createProblemMutation,
    getProblemsQuery,
  } from "@climblive/lib/queries";
  import { toastError, isDefined } from "@climblive/lib/utils";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let dialog: WaDialog | undefined = $state();
  let selectedContest: Contest | undefined = $state();
  let isCopying = $state(false);
  let copyProgress = $state(0);

  const contestsQuery = $derived(getAllContestsQuery());
  const selectedProblemsQuery = $derived(
    selectedContest ? getProblemsQuery(selectedContest.id) : undefined,
  );

  const availableContests = $derived.by(() => {
    if (!contestsQuery.data) return [];
    return contestsQuery.data
      .filter((c) => c.id !== contestId && !c.archived)
      .sort((a, b) => b.created.getTime() - a.created.getTime());
  });

  const problems = $derived(selectedProblemsQuery?.data);

  const handleOpen = () => {
    if (dialog) {
      dialog.open = true;
    }
  };

  const handleClose = () => {
    if (dialog) {
      dialog.open = false;
      selectedContest = undefined;
      copyProgress = 0;
    }
  };

  const handleSelectContest = (contest: Contest) => {
    selectedContest = contest;
  };

  const handleBack = () => {
    selectedContest = undefined;
  };

  const handleCopy = async () => {
    if (!problems || problems.length === 0) {
      toastError("No problems to copy");
      return;
    }

    isCopying = true;
    copyProgress = 0;
    let failCount = 0;

    const createMutation = createProblemMutation(contestId);

    for (let i = 0; i < problems.length; i++) {
      const problem = problems[i];
      const template: ProblemTemplate = {
        number: problem.number,
        holdColorPrimary: problem.holdColorPrimary,
        holdColorSecondary: problem.holdColorSecondary,
        description: problem.description,
        zone1Enabled: problem.zone1Enabled,
        zone2Enabled: problem.zone2Enabled,
        pointsZone1: problem.pointsZone1,
        pointsZone2: problem.pointsZone2,
        pointsTop: problem.pointsTop,
        flashBonus: problem.flashBonus,
      };

      await new Promise<void>((resolve) => {
        createMutation.mutate(template, {
          onSuccess: () => {
            resolve();
          },
          onError: () => {
            failCount++;
            resolve();
          },
        });
      });

      copyProgress = ((i + 1) / problems.length) * 100;
    }

    isCopying = false;

    if (failCount > 0) {
      toastError(
        `Failed to copy ${failCount} problem${failCount > 1 ? "s" : ""}`,
      );
    }

    handleClose();
  };

  const contestColumns: ColumnDefinition<Contest>[] = [
    {
      label: "Contest",
      mobile: true,
      render: renderContestName,
      width: "1fr",
    },
    {
      label: "Created",
      mobile: false,
      render: renderContestCreated,
      width: "max-content",
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
      render: renderPoints,
      width: "max-content",
    },
  ];
</script>

{#snippet renderContestName(contest: Contest)}
  {contest.name}
{/snippet}

{#snippet renderContestCreated(contest: Contest)}
  {contest.created.toLocaleDateString()}
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

{#snippet renderPoints({ pointsZone1, pointsZone2, pointsTop }: Problem)}
  {@const values = [pointsZone1, pointsZone2, pointsTop].filter(isDefined)}
  {values.join(" / ")} pts
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

<wa-button onclick={handleOpen} appearance="outlined" variant="neutral">
  Copy from another contest
  <wa-icon name="copy" slot="start"></wa-icon>
</wa-button>

<wa-dialog bind:this={dialog} label="Copy problems from another contest">
  {#if !selectedContest}
    <div class="dialog-content">
      <p>Select a contest to copy problems from:</p>
      {#if contestsQuery.isPending}
        <Loader />
      {:else if availableContests.length === 0}
        <p class="empty">No other contests available</p>
      {:else}
        <Table
          columns={contestColumns}
          data={availableContests}
          getId={({ id }) => id}
        ></Table>
      {/if}
    </div>
    <wa-button slot="footer" appearance="plain" onclick={handleClose}>
      Cancel
    </wa-button>
  {:else}
    <div class="dialog-content">
      {#if selectedProblemsQuery?.isPending}
        <Loader />
      {:else if problems && problems.length > 0}
        <p>
          {problems.length} problem{problems.length > 1 ? "s" : ""} from "{selectedContest.name}"
          will be copied:
        </p>
        <Table
          columns={problemColumns}
          data={problems}
          getId={({ id }) => id}
        />
        {#if isCopying}
          <wa-progress-bar value={copyProgress}></wa-progress-bar>
        {/if}
      {:else}
        <p class="empty">No problems found in this contest</p>
      {/if}
    </div>
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
</wa-dialog>

<style>
  .dialog-content {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
    min-height: 300px;
  }

  .number {
    display: flex;
    align-items: center;
    gap: var(--wa-space-s);
  }

  .empty {
    color: var(--wa-color-text-quiet);
    text-align: center;
    padding: var(--wa-space-xl);
  }

  wa-dialog {
    --width: 600px;
  }
</style>
