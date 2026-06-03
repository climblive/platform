<script lang="ts">
  import { type WaHideEvent } from "@awesome.me/webawesome";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/progress-bar/progress-bar.js";
  import "@awesome.me/webawesome/dist/components/tooltip/tooltip.js";
  import type {
    CompClassTemplate,
    CreateContendersArguments,
    ProblemTemplate,
  } from "@climblive/lib/models";
  import {
    createCompClassMutation,
    createContendersMutation,
    createProblemMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";

  const buttonId = $props.id();
  const primaryColors = [
    "#ef4444",
    "#f97316",
    "#f59e0b",
    "#eab308",
    "#84cc16",
    "#22c55e",
    "#10b981",
    "#14b8a6",
    "#06b6d4",
    "#0ea5e9",
    "#3b82f6",
    "#6366f1",
    "#8b5cf6",
    "#a855f7",
    "#d946ef",
    "#ec4899",
  ] as const;
  const totalProblems = 50;
  const totalSteps = 53;

  type Props = {
    contestId: number;
    disabled?: boolean;
    disabledReason?: string;
  };

  let { contestId, disabled = false, disabledReason }: Props = $props();

  let dialog: WaDialog | undefined = $state();
  let isRunning = $state(false);
  let completed = $state(false);
  let failed = $state(false);
  let progress = $state(0);
  let completedSteps = $state(0);
  let status = $state("Ready to create classes, problems, and tickets.");

  const createCompClass = $derived(createCompClassMutation(contestId));
  const createProblem = $derived(createProblemMutation(contestId));
  const createContenders = $derived(createContendersMutation(contestId));

  const openDialog = () => {
    if (disabled) {
      return;
    }

    completed = false;
    failed = false;
    progress = 0;
    completedSteps = 0;
    status = "Ready to create classes, problems, and tickets.";

    if (dialog) {
      dialog.open = true;
    }
  };

  const closeDialog = () => {
    if (isRunning) {
      return;
    }

    if (dialog) {
      dialog.open = false;
    }
  };

  const setProgress = (nextStatus: string, nextCompletedSteps: number) => {
    status = nextStatus;
    completedSteps = nextCompletedSteps;
    progress = (nextCompletedSteps / totalSteps) * 100;
  };

  const getProblemPointsTop = (index: number) => {
    return Math.round(25 + (275 * index) / (totalProblems - 1));
  };

  const getProblemTemplate = (index: number): ProblemTemplate => {
    const pointsTop = getProblemPointsTop(index);
    const hasFlashBonus = index % 4 === 0;
    const hasZone1 = index % 3 === 0;
    const hasZone2 = index % 6 === 0;
    const holdColorPrimary = primaryColors[index % primaryColors.length];
    const holdColorSecondary =
      index % 2 === 0
        ? primaryColors[(index + 5) % primaryColors.length]
        : undefined;

    return {
      number: index + 1,
      holdColorPrimary,
      holdColorSecondary,
      description: `Seeded problem ${index + 1}`,
      zone1Enabled: hasZone1,
      zone2Enabled: hasZone2,
      pointsZone1: hasZone1
        ? Math.max(1, Math.round(pointsTop * (hasZone2 ? 0.1 : 0.15)))
        : undefined,
      pointsZone2: hasZone2
        ? Math.max(1, Math.round(pointsTop * 0.2))
        : undefined,
      pointsTop,
      flashBonus: hasFlashBonus
        ? Math.max(1, Math.round(pointsTop * 0.05))
        : undefined,
    };
  };

  const getCompClasses = () => {
    const timeBegin = new Date();
    const timeEnd = new Date(timeBegin.getTime() + 12 * 60 * 60 * 1_000);

    return ["Males", "Females"].map<CompClassTemplate>((name) => ({
      name,
      timeBegin,
      timeEnd,
    }));
  };

  const handlePopulate = async () => {
    if (isRunning) {
      return;
    }

    isRunning = true;
    completed = false;
    failed = false;
    setProgress("Creating classes...", 0);

    try {
      const compClasses = getCompClasses();

      for (const [index, compClass] of compClasses.entries()) {
        await createCompClass.mutateAsync(compClass);
        setProgress(
          `Created class ${index + 1} of ${compClasses.length}.`,
          index + 1,
        );
      }

      for (let index = 0; index < totalProblems; index++) {
        await createProblem.mutateAsync(getProblemTemplate(index));
        setProgress(
          `Created problem ${index + 1} of ${totalProblems}.`,
          compClasses.length + index + 1,
        );
      }

      const args: CreateContendersArguments = {
        number: 200,
      };

      await createContenders.mutateAsync(args);
      setProgress("Created 200 tickets.", totalSteps);
      completed = true;
      status = "Contest populated with classes, problems, and tickets.";
    } catch {
      failed = true;
      status = "Contest population stopped before finishing.";
      toastError("Failed to populate contest.");
    } finally {
      isRunning = false;
    }
  };
</script>

<div class="actions">
  <wa-button
    id={buttonId}
    appearance="outlined"
    {disabled}
    onclick={openDialog}
  >
    Populate contest
    <wa-icon name="plus" slot="start"></wa-icon>
  </wa-button>
  {#if disabled && disabledReason}
    <wa-tooltip for={buttonId} placement="top-start"
      >{disabledReason}</wa-tooltip
    >
  {/if}
</div>

<wa-dialog
  bind:this={dialog}
  label="Populate contest"
  onwa-hide={(event: WaHideEvent) => {
    if (event.target !== dialog) {
      return;
    }

    if (isRunning) {
      event.preventDefault();
    }
  }}
>
  {#if isRunning || completed || failed}
    <div class="progress-content">
      <p>{status}</p>
      <wa-progress-bar value={progress}></wa-progress-bar>
      <small>{completedSteps} / {totalSteps} steps completed</small>
    </div>
  {:else}
    <div class="dialog-copy">
      <p>
        This will add 2 classes, 50 seeded problems, and 200 tickets to this
        contest.
      </p>
      <p>
        Both classes will start now and end in 12 hours. Problems will use
        varied colors, point values from 25 to 300, flash bonuses on some
        climbs, and zone scoring on some climbs.
      </p>
    </div>
  {/if}

  {#if completed || failed}
    <wa-button
      slot="footer"
      variant={completed ? "success" : "neutral"}
      onclick={closeDialog}
    >
      {#if completed}
        <wa-icon name="check" slot="start"></wa-icon>
      {/if}
      Close
    </wa-button>
  {:else}
    <wa-button
      slot="footer"
      appearance="plain"
      onclick={closeDialog}
      disabled={isRunning}
    >
      Cancel
    </wa-button>
    <wa-button
      slot="footer"
      variant="brand"
      onclick={handlePopulate}
      loading={isRunning}
    >
      Proceed
      <wa-icon name="plus" slot="start"></wa-icon>
    </wa-button>
  {/if}
</wa-dialog>

<style>
  wa-dialog::part(body) {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }

  .progress-content,
  .dialog-copy {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }

  p,
  small {
    margin: 0;
  }
</style>
