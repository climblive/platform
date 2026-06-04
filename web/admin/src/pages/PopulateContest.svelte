<script lang="ts" module>
  import { z } from "@climblive/lib/utils";

  const supportedClassNames = [
    "Males",
    "Females",
    "Seniors",
    "Kids",
    "Juniors",
  ] as const;

  const classNamesSchema = z
    .string()
    .transform((value, ctx) => {
      try {
        return JSON.parse(value);
      } catch {
        ctx.issues.push({
          code: "custom",
          input: value,
          message: "Select at least one class.",
        });

        return z.NEVER;
      }
    })
    .pipe(z.array(z.enum(supportedClassNames)).min(1));

  export const formSchema = z.object({
    contestLengthHours: z.coerce.number().int(),
    problemCount: z.coerce.number(),
    ticketCount: z.coerce.number(),
    problemMinPoints: z.coerce.number(),
    problemMaxPoints: z.coerce.number(),
    flashBonusPercentage: z.coerce.number(),
    zone1Percentage: z.coerce.number(),
    zone2Percentage: z.coerce.number(),
    classNames: classNamesSchema,
  });

  type PopulateContestFormData = z.infer<typeof formSchema>;
</script>

<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import "@awesome.me/webawesome/dist/components/progress-bar/progress-bar.js";
  import "@awesome.me/webawesome/dist/components/slider/slider.js";
  import "@awesome.me/webawesome/dist/components/tag/tag.js";
  import "@awesome.me/webawesome/dist/components/tooltip/tooltip.js";
  import { GenericForm, name, value } from "@climblive/lib/forms";
  import type {
    CompClassTemplate,
    CreateContendersArguments,
    ProblemTemplate,
  } from "@climblive/lib/models";
  import {
    createCompClassMutation,
    createContendersMutation,
    createProblemMutation,
    getContendersByContestQuery,
    getProblemsQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";

  const holdColors = [
    "#6f3601",
    "#dc3146",
    "#f46a45",
    "#fac22b",
    "#00ac49",
    "#2fbedc",
    "#0071ec",
    "#9951db",
    "#e66ba3",
    "#9194a2",
    "#000",
    "#fff",
  ] as const;

  type Props = {
    contestId: number;
  };

  let { contestId }: Props = $props();

  const contendersQuery = $derived(getContendersByContestQuery(contestId));
  const problemsQuery = $derived(getProblemsQuery(contestId));

  const contenders = $derived(contendersQuery.data);
  const problems = $derived(problemsQuery.data);

  let dialog = $state<WaDialog>();
  let isRunning = $state(false);
  let completed = $state(false);
  let failed = $state(false);
  let progress = $state(0);
  let completedSteps = $state(0);

  let problemCount = $derived(Math.min(50, 100 - (problems?.length ?? 0)));
  let ticketCount = $derived(Math.min(100, 500 - (contenders?.length ?? 0)));
  let contestLengthHours = $state(12);
  let problemMinPoints = $state(50);
  let problemMaxPoints = $state(300);
  let flashBonusPercentage = $state(5);
  let zone1Percentage = $state(15);
  let zone2Percentage = $state(20);
  let selectedClassNames = $state(supportedClassNames.slice(0, 2));

  const totalSteps = $derived(
    selectedClassNames.length + problemCount + (ticketCount > 0 ? 1 : 0),
  );

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
    selectedClassNames = [...supportedClassNames];

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

  const setProgress = (nextCompletedSteps: number) => {
    completedSteps = nextCompletedSteps;
    progress = totalSteps === 0 ? 100 : (nextCompletedSteps / totalSteps) * 100;
  };

  const getProblemPointsTop = (
    index: number,
    count: number,
    minPoints: number,
    maxPoints: number,
  ) => {
    if (count <= 1) {
      return maxPoints;
    }

    return Math.round(
      minPoints + ((maxPoints - minPoints) * index) / (count - 1),
    );
  };

  const getPointsByPercentage = (pointsTop: number, percentage: number) => {
    if (percentage <= 0) {
      return 0;
    }

    return Math.max(1, Math.round(pointsTop * (percentage / 100)));
  };

  const getProblemTemplate = (
    index: number,
    values: PopulateContestFormData,
  ): ProblemTemplate => {
    const pointsTop = getProblemPointsTop(
      index,
      values.problemCount,
      values.problemMinPoints,
      values.problemMaxPoints,
    );
    const hasFlashBonus = index % 4 === 0;
    const hasZone1 = index % 3 === 0;
    const hasZone2 = index % 6 === 0;
    const holdColorPrimary = holdColors[index % holdColors.length];

    return {
      number: index + 1,
      holdColorPrimary,
      description: `Seeded problem ${index + 1}`,
      zone1Enabled: hasZone1,
      zone2Enabled: hasZone2,
      pointsZone1: hasZone1
        ? getPointsByPercentage(pointsTop, values.zone1Percentage)
        : undefined,
      pointsZone2: hasZone2
        ? getPointsByPercentage(pointsTop, values.zone2Percentage)
        : undefined,
      pointsTop,
      flashBonus: hasFlashBonus
        ? getPointsByPercentage(pointsTop, values.flashBonusPercentage)
        : undefined,
    };
  };

  const getCompClasses = (values: PopulateContestFormData) => {
    const timeBegin = new Date();
    const timeEnd = new Date(
      timeBegin.getTime() + values.contestLengthHours * 60 * 60 * 1_000,
    );

    return values.classNames.map<CompClassTemplate>((name) => ({
      name,
      timeBegin,
      timeEnd,
    }));
  };

  const handleContestLengthChange = (event: Event) => {
    const target = event.target as HTMLInputElement;
    contestLengthHours = clamp(
      readInteger(target.value, contestLengthHours),
      1,
      72,
    );
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

  const handleRemoveClass = (name: (typeof supportedClassNames)[number]) => {
    selectedClassNames = selectedClassNames.filter(
      (className) => className !== name,
    );
  };

  const handleAddClass = (nextClass: (typeof supportedClassNames)[number]) => {
    selectedClassNames = [...selectedClassNames, nextClass];
  };

  const handlePopulate = async (values: PopulateContestFormData) => {
    if (isRunning) {
      return;
    }

    if (values.classNames.length === 0) {
      return;
    }

    isRunning = true;
    completed = false;
    failed = false;
    setProgress(0);

    try {
      const compClasses = getCompClasses(values);

      for (const [index, compClass] of compClasses.entries()) {
        await createCompClass.mutateAsync(compClass);
        setProgress(index + 1);
      }

      for (let index = 0; index < values.problemCount; index++) {
        await createProblem.mutateAsync(getProblemTemplate(index, values));
        setProgress(compClasses.length + index + 1);
      }

      if (values.ticketCount > 0) {
        const args: CreateContendersArguments = {
          number: values.ticketCount,
        };

        await createContenders.mutateAsync(args);
        setProgress(totalSteps);
      } else {
        setProgress(totalSteps);
      }

      completed = true;
    } catch {
      failed = true;
      toastError("Failed to populate contest.");
    } finally {
      isRunning = false;
    }
  };
</script>

<wa-button appearance="outlined" onclick={openDialog}>
  Populate with fake data
</wa-button>

<wa-dialog bind:this={dialog} open label="Populate contest">
  {#if isRunning || completed || failed}
    <wa-progress-bar value={progress}></wa-progress-bar>
    <small>{completedSteps} / {totalSteps} steps completed</small>
  {:else}
    <GenericForm schema={formSchema} submit={handlePopulate}>
      {@const nextClass = supportedClassNames.find(
        (name) => !selectedClassNames.includes(name),
      )}

      <input
        type="hidden"
        {@attach name("classNames")}
        value={JSON.stringify(selectedClassNames)}
      />
      <input
        type="hidden"
        {@attach name("problemMinPoints")}
        value={problemMinPoints}
      />
      <input
        type="hidden"
        {@attach name("problemMaxPoints")}
        value={problemMaxPoints}
      />

      <fieldset>
        <h4>Classes</h4>

        <div class="tags">
          {#each selectedClassNames as className (className)}
            <wa-tag
              size="s"
              variant="neutral"
              appearance="filled-outlined"
              with-remove
              onwa-remove={() => handleRemoveClass(className)}
            >
              {className}
            </wa-tag>
          {/each}

          {#if nextClass !== undefined}
            <wa-button
              type="button"
              appearance="plain"
              size="s"
              onclick={() => handleAddClass(nextClass)}
            >
              <wa-icon name="plus"></wa-icon>
            </wa-button>
          {/if}
        </div>

        <wa-number-input
          size="s"
          {@attach name("contestLengthHours")}
          label="Contest length"
          min="1"
          max="72"
          {@attach value(contestLengthHours.toString())}
          onchange={handleContestLengthChange}
        >
          <span slot="end">hours</span>
        </wa-number-input>

        <h4>Problems</h4>

        <div class="controls">
          <wa-number-input
            size="s"
            {@attach name("problemCount")}
            label="Number of problems"
            min="1"
            max={100 - (problems?.length ?? 0)}
            {@attach value(problemCount.toString())}
            onchange={handleProblemCountChange}
          ></wa-number-input>

          <wa-number-input
            size="s"
            {@attach name("ticketCount")}
            label="Number of tickets"
            min="0"
            max={500 - (contenders?.length ?? 0)}
            {@attach value(ticketCount.toString())}
            onchange={handleTicketCountChange}
          ></wa-number-input>

          <wa-number-input
            size="s"
            {@attach name("flashBonusPercentage")}
            label="Flash bonus"
            min="0"
            max="100"
            {@attach value(flashBonusPercentage.toString())}
            onchange={handleFlashBonusChange}
          >
            <span slot="end">%</span>
          </wa-number-input>

          <wa-number-input
            size="s"
            {@attach name("zone1Percentage")}
            label="Zone 1"
            min="0"
            max="100"
            {@attach value(zone1Percentage.toString())}
            onchange={handleZone1Change}
          >
            <span slot="end">%</span>
          </wa-number-input>

          <wa-number-input
            size="s"
            {@attach name("zone2Percentage")}
            label="Zone 2"
            min="0"
            max="100"
            {@attach value(zone2Percentage.toString())}
            onchange={handleZone2Change}
          >
            <span slot="end">%</span>
          </wa-number-input>
        </div>

        <wa-slider
          label="Point values"
          hint="Set the minimum and maximum top points."
          range
          min="0"
          max="1000"
          min-value={problemMinPoints}
          max-value={problemMaxPoints}
          step="25"
          with-tooltip
          oninput={handleProblemValueRangeChange}
        >
          <span slot="reference">0</span>
          <span slot="reference">500</span>
          <span slot="reference">1000</span>
        </wa-slider>

        <div class="footer-actions">
          <wa-button
            appearance="plain"
            type="button"
            onclick={closeDialog}
            disabled={isRunning}
          >
            Cancel
          </wa-button>
          <wa-button
            variant="neutral"
            type="submit"
            loading={isRunning}
            disabled={selectedClassNames.length === 0}
          >
            Proceed
          </wa-button>
        </div>
      </fieldset>
    </GenericForm>
  {/if}

  {#if completed || failed}
    <wa-button
      slot="footer"
      variant={completed ? "success" : "neutral"}
      onclick={closeDialog}
    >
      Close
    </wa-button>
  {/if}
</wa-dialog>

<style>
  wa-dialog::part(body) {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }

  .controls {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(12rem, 1fr));
    gap: var(--wa-space-s);
  }

  fieldset {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
    border: 0;
    padding: 0;
    margin: 0;
  }

  h4 {
    margin: 0;
  }

  .tags {
    display: flex;
    flex-wrap: wrap;
    gap: var(--wa-space-xs);
    align-items: center;
  }

  .footer-actions {
    display: flex;
    justify-content: flex-end;
    gap: var(--wa-space-xs);
  }

  small,
  span {
    margin: 0;
  }
</style>
