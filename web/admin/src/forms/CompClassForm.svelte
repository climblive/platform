<script lang="ts" module>
  import * as z from "zod";

  const twelveHours = 12 * 60 * 60 * 1_000;

  export const formSchema = z
    .object({
      name: z.string().min(1),
      description: z.string().optional(),
      timeBegin: z.coerce.date(),
      timeEnd: z.coerce.date(),
    })
    .superRefine((data, ctx) => {
      if (data.timeEnd.getTime() - data.timeBegin.getTime() > twelveHours) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "Total duration must not exceed 12 hours",
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
  import { GenericForm, name, value } from "@climblive/lib/forms";
  import type { CompClass } from "@climblive/lib/models";
  import "@shoelace-style/shoelace/dist/components/input/input.js";
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
</script>

<GenericForm {schema} {submit}>
  <fieldset>
    <sl-input
      size="small"
      use:name={"name"}
      label="Name"
      type="text"
      required
      use:value={data.name}
    ></sl-input>
    <sl-input
      size="small"
      use:name={"description"}
      label="Description"
      type="text"
      use:value={data.description}
    ></sl-input>
    <sl-input
      size="small"
      use:name={"timeBegin"}
      label="Start time"
      type="datetime-local"
      use:value={data.timeBegin
        ? format(data.timeBegin, "yyyy-MM-dd'T'HH:mm")
        : undefined}
    ></sl-input>
    <sl-input
      size="small"
      use:name={"timeEnd"}
      label="End time"
      type="datetime-local"
      use:value={data.timeEnd
        ? format(data.timeEnd, "yyyy-MM-dd'T'HH:mm")
        : undefined}
    ></sl-input>
    {@render children?.()}
  </fieldset>
</GenericForm>

<style>
  fieldset {
    display: flex;
    flex-direction: column;
    gap: var(--sl-spacing-small);
  }
</style>
