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
    .pipe(z.array(z.enum(supportedClassNames)));

  export const formSchema = z
    .object({
      contestLengthHours: z.coerce.number(),
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
  import { GenericForm, name } from "@climblive/lib/forms";
  import type {
    CompClassTemplate,
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
  import { add } from "date-fns";

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

  const { contestId }: Props = $props();

  const compClassesQuery = $derived(getCompClassesQuery(contestId));
  const contendersQuery = $derived(getContendersByContestQuery(contestId));
  const problemsQuery = $derived(getProblemsQuery(contestId));

  const compClasses = $derived(compClassesQuery.data);
  const contenders = $derived(contendersQuery.data);
  const problems = $derived(problemsQuery.data);

  let dialog = $state<WaDialog>();
  let populatorState = $state<"idle" | "pending" | "error" | "settled">("idle");

  const contestEmpty = $derived.by(() => {
    if (
      compClasses === undefined ||
      problems === undefined ||
      contenders === undefined
    ) {
      return false;
    }

    return (
      compClasses.length === 0 &&
      problems.length === 0 &&
      contenders.length === 0
    );
  });

  const maxExistingProblemNumber = $derived.by(() => {
    if (problems === undefined || problems.length === 0) {
      return 0;
    }

    return Math.max(...problems.map(({ number }) => number));
  });

  let totalSteps = $state(0);
  let completedSteps = $state(0);
  let progress = $derived(
    totalSteps === 0 ? 100 : (completedSteps / totalSteps) * 100,
  );

  const createCompClass = $derived(createCompClassMutation(contestId));
  const createProblem = $derived(createProblemMutation(contestId));
  const createContenders = $derived(createContendersMutation(contestId));

  const openDialog = () => {
    if (dialog) {
      dialog.open = true;
    }
  };

  const closeDialog = () => {
    if (dialog) {
      dialog.open = false;
    }
  };

  const incrementProgress = () => {
    completedSteps += 1;
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
    startNumberOffset: number,
    formData: PopulateContestFormData,
  ): ProblemTemplate => {
    const pointsTop = getProblemPointsTop(
      index,
      formData.problemCount,
      formData.problemMinPoints,
      formData.problemMaxPoints,
    );
    const hasFlashBonus = index % 4 === 0;
    const zone1Enabled = index % 3 === 0;
    const zone2Enabled = index % 6 === 0;
    const holdColorPrimary = holdColors[index % holdColors.length];
    const problemNumber = startNumberOffset + index;

    return {
      number: problemNumber,
      holdColorPrimary,
      zone1Enabled,
      zone2Enabled,
      pointsZone1: zone1Enabled
        ? getPointsByPercentage(pointsTop, formData.zone1Percentage)
        : undefined,
      pointsZone2: zone2Enabled
        ? getPointsByPercentage(pointsTop, formData.zone2Percentage)
        : undefined,
      pointsTop,
      flashBonus: hasFlashBonus
        ? getPointsByPercentage(pointsTop, formData.flashBonusPercentage)
        : undefined,
    };
  };

  const getCompClasses = (formData: PopulateContestFormData) => {
    const timeBegin = new Date();
    const timeEnd = add(timeBegin.getTime(), {
      hours: formData.contestLengthHours,
    });

    return formData.classNames.map<CompClassTemplate>((name) => ({
      name,
      timeBegin,
      timeEnd,
    }));
  };

  const handlePopulate = async (formData: PopulateContestFormData) => {
    populatorState = "pending";

    totalSteps =
      formData.classNames.length +
      formData.problemCount +
      (formData.ticketCount > 0 ? 1 : 0);

    const promises: Promise<unknown>[] = [];

    try {
      const compClasses = getCompClasses(formData);
      const firstProblemNumber = maxExistingProblemNumber + 1;

      for (const compClass of compClasses.values()) {
        promises.push(
          createCompClass.mutateAsync(compClass).then(incrementProgress),
        );
      }

      for (let index = 0; index < formData.problemCount; index++) {
        promises.push(
          createProblem
            .mutateAsync(
              getProblemTemplate(index, firstProblemNumber, formData),
            )
            .then(incrementProgress),
        );
      }

      if (formData.ticketCount > 0) {
        promises.push(
          createContenders
            .mutateAsync({
              number: formData.ticketCount,
            })
            .then(incrementProgress),
        );
      }

      await Promise.all(promises);
      populatorState = "settled";
    } catch {
      populatorState = "error";
      toastError("Failed to populate contest.");
    }
  };
</script>

<wa-button appearance="outlined" onclick={openDialog} disabled={!contestEmpty}>
  Populate with fake data
</wa-button>

<wa-dialog bind:this={dialog} label="Populate contest">
  {#if populatorState !== "idle"}
    <wa-progress-bar value={progress}></wa-progress-bar>
    <wa-button
      slot="footer"
      variant={progress === 100 ? "success" : "neutral"}
      onclick={closeDialog}
    >
      Close
    </wa-button>
  {:else}
    <GenericForm schema={formSchema} submit={handlePopulate}>
      <div class="content">
        <h4>Classes</h4>

        <wa-select
          size="s"
          multiple
          with-clear
          {@attach name("className")}
          label="Classes"
          hint="Select the classes to create."
          required
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
          value="12"
          required
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
          value="50"
          required
        ></wa-number-input>

        <wa-slider
          {@attach name("problemPoints")}
          label="Point values"
          hint="Set the minimum and maximum top points."
          range
          min="0"
          max="1000"
          min-value="50"
          max-value="300"
          step="25"
          with-tooltip
          required
          size="s"
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
          value="5"
          required
        >
          <span slot="end">%</span>
        </wa-number-input>

        <wa-number-input
          size="s"
          {@attach name("zone1Percentage")}
          label="Zone 1"
          min="0"
          max="100"
          value="15"
          required
        >
          <span slot="end">%</span>
        </wa-number-input>

        <wa-number-input
          size="s"
          {@attach name("zone2Percentage")}
          label="Zone 2"
          min="0"
          max="100"
          value="25"
          required
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
          value="100"
          required
        ></wa-number-input>

        <div class="footer-actions">
          <wa-button appearance="plain" type="button" onclick={closeDialog}>
            Cancel
          </wa-button>
          <wa-button variant="neutral" type="submit"> Proceed </wa-button>
        </div>
      </div>
    </GenericForm>
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
</style>
