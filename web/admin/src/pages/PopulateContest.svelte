<script lang="ts" module>
  import { z } from "@climblive/lib/utils";

  const supportedClassNames = [
    "Males",
    "Females",
    "Boys",
    "Girls",
    "Seniors",
  ] as const;

  const classNamesSchema = z
    .union([z.enum(supportedClassNames), z.array(z.enum(supportedClassNames))])
    .transform((value) => {
      if (Array.isArray(value)) {
        return value;
      }

      return [value];
    })
    .pipe(z.array(z.enum(supportedClassNames)).min(1));

  export const formSchema = z
    .object({
      contestLengthHours: z.coerce.number().int(),
      problemCount: z.coerce.number(),
      ticketCount: z.coerce.number(),
      problemPoints: z
        .union([z.coerce.number(), z.array(z.coerce.number())])
        .transform((value) => {
          return Array.isArray(value) ? value : [value];
        }),
      flashBonusPercentage: z.coerce.number(),
      zone1Percentage: z.coerce.number(),
      zone2Percentage: z.coerce.number(),
      className: classNamesSchema,
    })
    .transform((value) => ({
      ...value,
      problemMinPoints: Math.min(...value.problemPoints),
      problemMaxPoints: Math.max(...value.problemPoints),
      classNames: value.className,
    }));

  type PopulateContestFormData = z.infer<typeof formSchema>;
</script>

<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/progress-bar/progress-bar.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import "@awesome.me/webawesome/dist/components/slider/slider.js";
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
    getCompClassesQuery,
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

  const defaultContestLengthHours = 12;
  const defaultProblemMinPoints = 50;
  const defaultProblemMaxPoints = 300;
  const defaultFlashBonusPercentage = 5;
  const defaultZone1Percentage = 15;
  const defaultZone2Percentage = 20;

  let { contestId }: Props = $props();

  const compClassesQuery = $derived(getCompClassesQuery(contestId));
  const contendersQuery = $derived(getContendersByContestQuery(contestId));
  const problemsQuery = $derived(getProblemsQuery(contestId));

  const compClasses = $derived(compClassesQuery.data);
  const contenders = $derived(contendersQuery.data);
  const problems = $derived(problemsQuery.data);

  let dialog = $state<WaDialog>();
  let isRunning = $state(false);
  let completed = $state(false);
  let failed = $state(false);
  let progress = $state(0);
  let completedSteps = $state(0);
  let submittedValues = $state<PopulateContestFormData>();
  let formVersion = $state(0);

  const defaultProblemCount = $derived(
    Math.min(50, 100 - (problems?.length ?? 0)),
  );

  const defaultTicketCount = $derived(
    Math.min(100, 500 - (contenders?.length ?? 0)),
  );

  const contestNotEmpty = $derived.by(() => {
    if (
      compClasses === undefined ||
      problems === undefined ||
      contenders === undefined
    ) {
      return true;
    }

    return (
      compClasses.length > 0 || problems.length > 0 || contenders.length > 0
    );
  });

  const maxExistingProblemNumber = $derived.by(() => {
    if (problems === undefined || problems.length === 0) {
      return 0;
    }

    return Math.max(...problems.map(({ number }) => number));
  });

  const totalSteps = $derived(
    submittedValues === undefined
      ? 0
      : submittedValues.classNames.length +
          submittedValues.problemCount +
          (submittedValues.ticketCount > 0 ? 1 : 0),
  );

  const createCompClass = $derived(createCompClassMutation(contestId));
  const createProblem = $derived(createProblemMutation(contestId));
  const createContenders = $derived(createContendersMutation(contestId));

  const openDialog = () => {
    if (contestNotEmpty) {
      return;
    }

    completed = false;
    failed = false;
    progress = 0;
    completedSteps = 0;
    submittedValues = undefined;
    formVersion += 1;

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
    startNumber: number,
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
    const problemNumber = startNumber + index;

    return {
      number: problemNumber,
      holdColorPrimary,
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
    submittedValues = values;
    setProgress(0);

    try {
      const compClasses = getCompClasses(values);
      const firstProblemNumber = maxExistingProblemNumber + 1;

      for (const [index, compClass] of compClasses.entries()) {
        await createCompClass.mutateAsync(compClass);
        setProgress(index + 1);
      }

      for (let index = 0; index < values.problemCount; index++) {
        await createProblem.mutateAsync(
          getProblemTemplate(index, firstProblemNumber, values),
        );
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

<wa-button
  appearance="outlined"
  onclick={openDialog}
  disabled={contestNotEmpty}
>
  Populate with fake data
</wa-button>

<wa-dialog bind:this={dialog} label="Populate contest">
  {#if isRunning || completed || failed}
    <wa-progress-bar value={progress}></wa-progress-bar>
    <small>{completedSteps} / {totalSteps} steps completed</small>
  {:else}
    {#key formVersion}
      <GenericForm schema={formSchema} submit={handlePopulate}>
        <div class="content">
          <h4>Classes</h4>

          <wa-select
            size="s"
            multiple
            with-clear
            max-options-visible="3"
            {@attach name("className")}
            label="Classes"
            hint="Select the classes to create."
          >
            {#each supportedClassNames as className, index (className)}
              <wa-option value={className} selected={index < 2}>
                {className}
              </wa-option>
            {/each}
          </wa-select>

          <wa-number-input
            size="s"
            {@attach name("contestLengthHours")}
            label="Contest length"
            min="1"
            max="72"
            {@attach value(defaultContestLengthHours.toString())}
          >
            <span slot="end">hours</span>
          </wa-number-input>

          <h4>Problems</h4>

          <wa-number-input
            size="s"
            {@attach name("problemCount")}
            label="Number of problems"
            min="0"
            max="100"
            {@attach value(defaultProblemCount.toString())}
          ></wa-number-input>

          <wa-slider
            {@attach name("problemPoints")}
            label="Point values"
            hint="Set the minimum and maximum top points."
            range
            min="0"
            max="1000"
            min-value={defaultProblemMinPoints}
            max-value={defaultProblemMaxPoints}
            step="25"
            with-tooltip
          >
            <span slot="reference">0</span>
            <span slot="reference">500</span>
            <span slot="reference">1000</span>
          </wa-slider>

          <wa-number-input
            size="s"
            {@attach name("flashBonusPercentage")}
            label="Flash bonus"
            min="0"
            max="100"
            {@attach value(defaultFlashBonusPercentage.toString())}
          >
            <span slot="end">%</span>
          </wa-number-input>

          <wa-number-input
            size="s"
            {@attach name("zone1Percentage")}
            label="Zone 1"
            min="0"
            max="100"
            {@attach value(defaultZone1Percentage.toString())}
          >
            <span slot="end">%</span>
          </wa-number-input>

          <wa-number-input
            size="s"
            {@attach name("zone2Percentage")}
            label="Zone 2"
            min="0"
            max="100"
            {@attach value(defaultZone2Percentage.toString())}
          >
            <span slot="end">%</span>
          </wa-number-input>

          <h4>Tickets</h4>

          <wa-number-input
            size="s"
            {@attach name("ticketCount")}
            label="Number of tickets"
            min="0"
            max="500"
            {@attach value(defaultTicketCount.toString())}
          ></wa-number-input>

          <div class="footer-actions">
            <wa-button
              appearance="plain"
              type="button"
              onclick={closeDialog}
              disabled={isRunning}
            >
              Cancel
            </wa-button>
            <wa-button variant="neutral" type="submit" loading={isRunning}>
              Proceed
            </wa-button>
          </div>
        </div>
      </GenericForm>
    {/key}
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
  .content {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
    margin: 0;
  }

  h4 {
    margin: 0;

    &:not(:first-of-type) {
      margin-block-start: var(--wa-space-m);
    }
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
