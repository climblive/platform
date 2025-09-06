<script lang="ts" module>
  import * as z from "zod";

  const oneMonth = 31 * 24 * 60 * 60 * 1_000;

  export const formSchema = z
    .object({
      name: z.string().min(1),
      description: z.string().optional(),
      timeBegin: z.coerce.date(),
      timeEnd: z.coerce.date(),
    })
    .superRefine((data, ctx) => {
      if (data.timeEnd.getTime() - data.timeBegin.getTime() > oneMonth) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "Total duration must not exceed 31 days",
          path: ["timeEnd"],
        });
      }

      if (data.timeEnd <= data.timeBegin) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "Time must follow chronological order",
          path: ["timeEnd"],
        });
      }
    });
</script>

<script lang="ts">
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import type WaInput from "@awesome.me/webawesome/dist/components/input/input.js";
  import { GenericForm, name } from "@climblive/lib/forms";
  import type { CompClass } from "@climblive/lib/models";
  import { format } from "date-fns";
  import { type Snippet } from "svelte";

  type T = $$Generic<Partial<CompClass>>;

  interface Props {
    data: Partial<T>;
    schema: z.ZodType<T, z.ZodTypeDef, T>;
    submit: (value: T) => void;
    children?: Snippet;
  }

  let { data, schema, submit, children }: Props = $props();

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

      const diff = begin.getTime() - previousTimeBegin.getTime();
      console.log("diff", diff);
      if (diff > 0) {
        end.setTime(end.getTime() + diff);
        timeEndInput.value = format(end, "yyyy-MM-dd'T'HH:mm");
      }
    } finally {
      previousTimeBegin = begin;
    }
  }
</script>

<GenericForm {schema} {submit}>
  <fieldset>
    <wa-input
      size="small"
      {@attach name("name")}
      label="Name"
      type="text"
      required
      value={data.name}
      placeholder="Males or Females"
    ></wa-input>
    <wa-input
      size="small"
      {@attach name("description")}
      label="Description"
      type="text"
      value={data.description}
    ></wa-input>
    <wa-input
      bind:this={timeBeginInput}
      size="small"
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
      size="small"
      {@attach name("timeEnd")}
      label="End time"
      type="datetime-local"
      value={data.timeEnd
        ? format(data.timeEnd, "yyyy-MM-dd'T'HH:mm")
        : undefined}
    ></wa-input>
    {@render children?.()}
  </fieldset>
</GenericForm>

<style>
  fieldset {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }
</style>
