<script lang="ts">
  import { type WaHideEvent } from "@awesome.me/webawesome";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import "@awesome.me/webawesome/dist/components/progress-bar/progress-bar.js";
  import "@awesome.me/webawesome/dist/components/slider/slider.js";
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

  type Props = {
    contestId: number;
  };

  let { contestId }: Props = $props();

  let dialog: WaDialog | undefined = $state();
  let isRunning = $state(false);
  let completed = $state(false);
  let failed = $state(false);
  let progress = $state(0);
  let completedSteps = $state(0);
  let status = $state("Ready to create classes, problems, and tickets.");

  let problemCount = $state(50);
  let ticketCount = $state(200);
  let problemMinPoints = $state(25);
  let problemMaxPoints = $state(300);
  let flashBonusPercentage = $state(5);
  let zone1Percentage = $state(15);
  let zone2Percentage = $state(20);

  const totalSteps = $derived(2 + problemCount + (ticketCount > 0 ? 1 : 0));

  const summary = $derived.by(() => {
    return `This will add 2 classes, ${problemCount} seeded problem${
      problemCount === 1 ? "" : "s"
    }, and ${ticketCount} ticket${ticketCount === 1 ? "" : "s"} to this contest.`;
  });

  const scoringSummary = $derived.by(() => {
    return `Problems will range from ${problemMinPoints} to ${problemMaxPoints} points, with a ${flashBonusPercentage}% flash bonus on some climbs, ${zone1Percentage}% zone 1 points on some climbs, and ${zone2Percentage}% zone 2 points on some climbs.`;
  });

  const createCompClass = $derived(createCompClassMutation(contestId));
  const createProblem = $derived(createProblemMutation(contestId));
  const createContenders = $derived(createContendersMutation(contestId));

  const clamp = (value: number, min: number, max: number) => {
    return Math.min(max, Math.max(min, value));
  };

  const readInteger = (value: string | number, fallback: number) => {
    const parsed = Number(value);

    if (Number.isNaN(parsed)) {
      return fallback;
    }

    return Math.round(parsed);
  };

  const openDialog = () => {
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
    progress = totalSteps === 0 ? 100 : (nextCompletedSteps / totalSteps) * 100;
  };

  const getProblemPointsTop = (index: number) => {
    if (problemCount <= 1) {
      return problemMaxPoints;
    }

    return Math.round(
      problemMinPoints +
        ((problemMaxPoints - problemMinPoints) * index) / (problemCount - 1),
    );
  };

  const getPointsByPercentage = (pointsTop: number, percentage: number) => {
    if (percentage <= 0) {
      return 0;
    }

    return Math.max(1, Math.round(pointsTop * (percentage / 100)));
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
        ? getPointsByPercentage(pointsTop, zone1Percentage)
        : undefined,
      pointsZone2: hasZone2
        ? getPointsByPercentage(pointsTop, zone2Percentage)
        : undefined,
      pointsTop,
      flashBonus: hasFlashBonus
        ? getPointsByPercentage(pointsTop, flashBonusPercentage)
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

  const handleProblemCountChange = (event: Event) => {
    const target = event.target as HTMLInputElement;
    problemCount = clamp(readInteger(target.value, problemCount), 1, 100);
  };

  const handleTicketCountChange = (event: Event) => {
    const target = event.target as HTMLInputElement;
    ticketCount = clamp(readInteger(target.value, ticketCount), 0, 500);
  };

  const handleFlashBonusChange = (event: Event) => {
    const target = event.target as HTMLInputElement;
    flashBonusPercentage = clamp(
      readInteger(target.value, flashBonusPercentage),
      0,
      100,
    );
  };

  const handleZone1Change = (event: Event) => {
    const target = event.target as HTMLInputElement;
    zone1Percentage = clamp(readInteger(target.value, zone1Percentage), 0, 100);
  };

  const handleZone2Change = (event: Event) => {
    const target = event.target as HTMLInputElement;
    zone2Percentage = clamp(readInteger(target.value, zone2Percentage), 0, 100);
  };

  const handleProblemValueRangeChange = (event: Event) => {
    const target = event.target as HTMLElement & {
      minValue: number;
      maxValue: number;
    };

    problemMinPoints = clamp(Math.round(target.minValue), 0, 1000);
    problemMaxPoints = clamp(
      Math.round(target.maxValue),
      problemMinPoints,
      1000,
    );
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

      for (let index = 0; index < problemCount; index++) {
        await createProblem.mutateAsync(getProblemTemplate(index));
        setProgress(
          `Created problem ${index + 1} of ${problemCount}.`,
          compClasses.length + index + 1,
        );
      }

      if (ticketCount > 0) {
        const args: CreateContendersArguments = {
          number: ticketCount,
        };

        await createContenders.mutateAsync(args);
        setProgress(`Created ${ticketCount} tickets.`, totalSteps);
      } else {
        setProgress("Skipped ticket creation.", totalSteps);
      }

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
  <wa-button appearance="outlined" onclick={openDialog}>
    Populate contest
    <wa-icon name="plus" slot="start"></wa-icon>
  </wa-button>
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
      <p>{summary}</p>
      <p>
        Both classes will start now and end in 12 hours. {scoringSummary}
      </p>

      <div class="controls">
        <wa-number-input
          size="s"
          label="Number of problems"
          min="1"
          max="100"
          value={problemCount.toString()}
          onchange={handleProblemCountChange}
        ></wa-number-input>

        <wa-number-input
          size="s"
          label="Number of tickets"
          min="0"
          max="500"
          value={ticketCount.toString()}
          onchange={handleTicketCountChange}
        ></wa-number-input>

        <wa-number-input
          size="s"
          label="Flash bonus (%)"
          min="0"
          max="100"
          value={flashBonusPercentage.toString()}
          onchange={handleFlashBonusChange}
        ></wa-number-input>

        <wa-number-input
          size="s"
          label="Zone 1 (%)"
          min="0"
          max="100"
          value={zone1Percentage.toString()}
          onchange={handleZone1Change}
        ></wa-number-input>

        <wa-number-input
          size="s"
          label="Zone 2 (%)"
          min="0"
          max="100"
          value={zone2Percentage.toString()}
          onchange={handleZone2Change}
        ></wa-number-input>
      </div>

      <div class="range">
        <div class="range-header">
          <span>Problem value range</span>
          <small>{problemMinPoints} - {problemMaxPoints} pts</small>
        </div>

        <wa-slider
          label="Problem value range"
          hint="Set the minimum and maximum top points."
          range
          min="0"
          max="1000"
          min-value={problemMinPoints}
          max-value={problemMaxPoints}
          step="1"
          with-tooltip
          oninput={handleProblemValueRangeChange}
        >
          <span slot="reference">0</span>
          <span slot="reference">500</span>
          <span slot="reference">1000</span>
        </wa-slider>
      </div>
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
  .dialog-copy,
  .range {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }

  .controls {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(12rem, 1fr));
    gap: var(--wa-space-s);
  }

  .range-header {
    display: flex;
    justify-content: space-between;
    gap: var(--wa-space-s);
    align-items: baseline;
  }

  p,
  small,
  span {
    margin: 0;
  }
</style>
