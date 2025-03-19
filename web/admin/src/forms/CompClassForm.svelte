<script lang="ts">
  import { GenericForm, name, value } from "@climblive/lib/forms";
  import { type CompClassTemplate } from "@climblive/lib/models";
  import "@shoelace-style/shoelace/dist/components/input/input.js";
  import { type Snippet } from "svelte";
  import * as z from "zod";

  const formSchema: z.ZodType<CompClassTemplate> = z.object({
    name: z.string().min(1),
    description: z.string().optional(),
    timeBegin: z.coerce.date(),
    timeEnd: z.coerce.date(),
  });

  interface Props {
    data: Partial<CompClassTemplate>;
    submit: (patch: CompClassTemplate) => void;
    children?: Snippet;
  }

  let { data, submit, children }: Props = $props();
</script>

<GenericForm schema={formSchema} {submit}>
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
      label="Time begin"
      type="datetime-local"
      use:value={data.timeBegin?.toISOString()}
    ></sl-input>
    <sl-input
      size="small"
      use:name={"timeEnd"}
      label="Time end"
      type="datetime-local"
      use:value={data.timeEnd?.toISOString()}
    ></sl-input>
    {@render children?.()}
  </fieldset>
</GenericForm>

<style>
  fieldset {
    display: flex;
    flex-direction: column;
    gap: var(--sl-spacing-small);
    padding: var(--sl-spacing-medium);
  }
</style>
