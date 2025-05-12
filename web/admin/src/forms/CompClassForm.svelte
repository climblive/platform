<script lang="ts">
  import { GenericForm, name, value } from "@climblive/lib/forms";
  import type { CompClass } from "@climblive/lib/models";
  import "@shoelace-style/shoelace/dist/components/input/input.js";
  import { format } from "date-fns";
  import { type Snippet } from "svelte";
  import * as z from "zod";

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
