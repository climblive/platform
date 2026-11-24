<script lang="ts" module>
  import { categoryEligibilitySchema } from "@climblive/lib/models";
  import { z } from "@climblive/lib/utils";

  const oneMonth = 31 * 24 * 60 * 60 * 1_000;

  const optionalBound = (minimum: number, maximum: number) =>
    z.preprocess(
      (value) =>
        value === undefined ||
        value === null ||
        (typeof value === "string" && value.trim() === "")
          ? null
          : value,
      z.coerce.number().int().min(minimum).max(maximum).nullable(),
    );

  export const formSchema = z
    .object({
      name: z.string().min(1),
      eligibilityBasis: z.enum(["age", "birth_year"]),
      minAge: optionalBound(0, 2 ** 31 - 1),
      maxAge: optionalBound(0, 2 ** 31 - 1),
      minBirthYear: optionalBound(1, 9999),
      maxBirthYear: optionalBound(1, 9999),
      timeBegin: z.coerce.date(),
      timeEnd: z.coerce.date(),
    })
    .pipe(
      categoryEligibilitySchema.safeExtend({
        name: z.string().min(1),
        timeBegin: z.date(),
        timeEnd: z.date(),
      }),
    )
    .check((ctx) => {
      const data = ctx.value;

      if (data.timeEnd.getTime() - data.timeBegin.getTime() > oneMonth) {
        ctx.issues.push({
          code: "custom",
          message: "Total duration must not exceed 31 days",
          input: data,
          path: ["timeEnd"],
        });
      }

      if (data.timeEnd <= data.timeBegin) {
        ctx.issues.push({
          code: "custom",
          message: "Time must follow chronological order",
          input: data,
          path: ["timeEnd"],
        });
      }
    });
</script>

<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import type WaNumberInput from "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import "@awesome.me/webawesome/dist/components/radio-group/radio-group.js";
  import type WaRadioGroup from "@awesome.me/webawesome/dist/components/radio-group/radio-group.js";
  import "@awesome.me/webawesome/dist/components/radio/radio.js";
  import type WaInput from "@awesome.me/webawesome/dist/components/input/input.js";
  import { GenericForm, name } from "@climblive/lib/forms";
  import type { Category } from "@climblive/lib/models";
  import { formatCategoryEligibility } from "@climblive/lib/utils";
  import { addSeconds, format } from "date-fns";
  import { type Snippet } from "svelte";

  type T = $$Generic<Partial<Category>>;

  interface Props {
    data: Partial<Category>;
    schema: z.ZodType<T, unknown>;
    submit: (value: T) => void;
    children?: Snippet;
  }

  let { data, schema, submit, children }: Props = $props();

  let eligibilityBasis = $derived(data.eligibilityBasis ?? "age");
  let minAge = $derived(data.minAge?.toString() ?? "");
  let maxAge = $derived(data.maxAge?.toString() ?? "");
  let minBirthYear = $derived(data.minBirthYear?.toString() ?? "");
  let maxBirthYear = $derived(data.maxBirthYear?.toString() ?? "");
  const eligibilityHintId = $props.id();

  const eligibilityPreview = $derived(
    categoryEligibilitySchema.safeParse({
      eligibilityBasis,
      minAge:
        eligibilityBasis === "age" && minAge !== "" ? Number(minAge) : null,
      maxAge:
        eligibilityBasis === "age" && maxAge !== "" ? Number(maxAge) : null,
      minBirthYear:
        eligibilityBasis === "birth_year" && minBirthYear !== ""
          ? Number(minBirthYear)
          : null,
      maxBirthYear:
        eligibilityBasis === "birth_year" && maxBirthYear !== ""
          ? Number(maxBirthYear)
          : null,
    }),
  );

  let previousTimeBegin = $derived(data.timeBegin);

  let timeBeginInput: WaInput | undefined = $state();
  let timeEndInput: WaInput | undefined = $state();

  function handleTimeBeginChange() {
    let begin: Date;
    let end: Date;

    if (timeBeginInput && timeBeginInput.value) {
      begin = new Date(timeBeginInput.value);
    } else {
      return;
    }

    try {
      if (!timeEndInput || !timeEndInput.value || !previousTimeBegin) {
        return;
      }

      end = new Date(timeEndInput.value);
      const diff = end.getTime() - previousTimeBegin.getTime();

      end.setTime(begin.getTime() + diff);
      timeEndInput.value = format(end, "yyyy-MM-dd'T'HH:mm");
    } finally {
      previousTimeBegin = begin;
    }
  }

  const extendTimeEnd = (minutes: number) => {
    if (!timeEndInput || !timeEndInput.value) {
      return;
    }

    const end = new Date(timeEndInput.value);
    timeEndInput.value = format(
      addSeconds(end, minutes * 60),
      "yyyy-MM-dd'T'HH:mm",
    );
  };
</script>

<GenericForm {schema} {submit}>
  <fieldset>
    <wa-input
      size="s"
      {@attach name("name")}
      label="Name"
      type="text"
      required
      value={data.name}
      placeholder="Males or Females"
    ></wa-input>
    <div class="eligibility">
      <wa-radio-group
        size="s"
        label="Define category by"
        orientation="horizontal"
        {@attach name("eligibilityBasis")}
        value={eligibilityBasis}
        required
        onchange={(event: Event) => {
          eligibilityBasis = String(
            (event.currentTarget as WaRadioGroup).value,
          );
        }}
      >
        <wa-radio appearance="button" value="age">Age range</wa-radio>
        <wa-radio appearance="button" value="birth_year">Birth year</wa-radio>
      </wa-radio-group>
      <div class="age-span">
        {#if eligibilityBasis === "age"}
          <wa-number-input
            size="s"
            {@attach name("minAge")}
            label="Minimum age"
            placeholder="Any"
            min="0"
            max={2 ** 31 - 1}
            step="1"
            value={minAge}
            aria-describedby={eligibilityHintId}
            oninput={(event: Event) => {
              minAge = (event.currentTarget as WaNumberInput).value ?? "";
            }}
          ></wa-number-input>
          <span aria-hidden="true">-</span>
          <wa-number-input
            size="s"
            {@attach name("maxAge")}
            label="Maximum age"
            placeholder="Any"
            min="0"
            max={2 ** 31 - 1}
            step="1"
            value={maxAge}
            aria-describedby={eligibilityHintId}
            oninput={(event: Event) => {
              maxAge = (event.currentTarget as WaNumberInput).value ?? "";
            }}
          ></wa-number-input>
        {:else}
          <wa-number-input
            size="s"
            {@attach name("minBirthYear")}
            label="Earliest birth year"
            placeholder="Any"
            min="1"
            max="9999"
            step="1"
            value={minBirthYear}
            aria-describedby={eligibilityHintId}
            oninput={(event: Event) => {
              minBirthYear = (event.currentTarget as WaNumberInput).value ?? "";
            }}
          ></wa-number-input>
          <span aria-hidden="true">-</span>
          <wa-number-input
            size="s"
            {@attach name("maxBirthYear")}
            label="Latest birth year"
            placeholder="Any"
            min="1"
            max="9999"
            step="1"
            value={maxBirthYear}
            aria-describedby={eligibilityHintId}
            oninput={(event: Event) => {
              maxBirthYear = (event.currentTarget as WaNumberInput).value ?? "";
            }}
          ></wa-number-input>
        {/if}
      </div>
      <p id={eligibilityHintId} class="hint">
        {#if eligibilityBasis === "age"}
          Age is measured at competition start.
        {/if}
        Leave either end empty for no limit.
      </p>
      <p class="summary" aria-live="polite">
        {#if eligibilityPreview.success}
          {formatCategoryEligibility(eligibilityPreview.data)}
        {/if}
      </p>
    </div>
    <wa-input
      bind:this={timeBeginInput}
      size="s"
      onchange={handleTimeBeginChange}
      {@attach name("timeBegin")}
      label="Start time"
      type="datetime-local"
      value={data.timeBegin
        ? format(data.timeBegin, "yyyy-MM-dd'T'HH:mm")
        : undefined}
    ></wa-input>
    <wa-input
      bind:this={timeEndInput}
      size="s"
      {@attach name("timeEnd")}
      label="End time"
      type="datetime-local"
      value={data.timeEnd
        ? format(data.timeEnd, "yyyy-MM-dd'T'HH:mm")
        : undefined}
    ></wa-input>
    <div class="quick-controls">
      <wa-button pill size="s" onclick={() => extendTimeEnd(5)}
        >+5 min</wa-button
      >
      <wa-button pill size="s" onclick={() => extendTimeEnd(15)}
        >+15 min</wa-button
      >
    </div>
    {@render children?.()}
  </fieldset>
</GenericForm>

<style>
  fieldset {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }

  .quick-controls {
    display: flex;
    gap: var(--wa-space-xs);
  }

  .age-span {
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
    gap: var(--wa-space-s);
    align-items: end;

    wa-number-input {
      min-width: 0;
    }

    span {
      padding-block: var(--wa-space-s);
    }
  }

  .eligibility {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }

  .hint {
    margin: 0;
    color: var(--wa-color-text-quiet);
    font-size: var(--wa-font-size-s);
  }

  .summary {
    margin: 0;
    min-height: 1lh;
  }
</style>
