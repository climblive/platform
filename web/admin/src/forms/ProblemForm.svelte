<script lang="ts">
  import { GenericForm, name, value } from "@climblive/lib/forms";
  import { type ProblemTemplate } from "@climblive/lib/models";
  import "@shoelace-style/shoelace/dist/components/input/input.js";
  import { type Snippet } from "svelte";
  import * as z from "zod";

  const formSchema: z.ZodType<ProblemTemplate> = z.object({
    number: z.coerce.number(),
    holdColorPrimary: z.string(),
    holdColorSecondary: z.string().optional(),
    name: z.string().optional(),
    description: z.string().optional(),
    pointsTop: z.coerce.number(),
    pointsZone: z.coerce.number(),
    flashBonus: z.coerce.number().optional(),
  });

  interface Props {
    data: Partial<ProblemTemplate>;
    submit: (value: ProblemTemplate) => void;
    children?: Snippet;
  }

  let { data, submit, children }: Props = $props();
</script>

<GenericForm schema={formSchema} {submit}>
  <fieldset>
    <sl-input
      size="small"
      use:name={"number"}
      label="Number"
      type="number"
      required
      use:value={data.number}
    ></sl-input>
    <sl-input
      size="small"
      use:name={"holdColorPrimary"}
      label="Primary hold color"
      type="text"
      required
      use:value={data.holdColorPrimary}
    ></sl-input>
    <sl-input
      size="small"
      use:name={"holdColorSecondary"}
      label="Secondary hold color"
      type="text"
      use:value={data.holdColorSecondary}
    ></sl-input>
    <sl-input
      size="small"
      use:name={"name"}
      label="Name"
      type="text"
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
      use:name={"pointsTop"}
      label="Points for top"
      type="number"
      required
      use:value={data.pointsTop}
    ></sl-input>
    <sl-input
      size="small"
      use:name={"pointsZone"}
      label="Points for zone"
      type="number"
      required
      use:value={data.pointsZone}
    ></sl-input>
    <sl-input
      size="small"
      use:name={"flashBonus"}
      label="Flash bonus"
      type="number"
      use:value={data.flashBonus}
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
