<script lang="ts" module>
  import * as z from "zod";

  export const formSchema = z.object({
    number: z.coerce.number(),
    holdColorPrimary: z.string().regex(/^#([0-9a-fA-F]{3}){1,2}$/),
    holdColorSecondary: z.string().optional(),
    description: z.string().optional(),
    pointsTop: z.coerce.number(),
    pointsZone: z.coerce.number(),
    flashBonus: z.coerce.number().optional(),
  });
</script>

<script lang="ts">
  import "@awesome.me/webawesome/dist/components/color-picker/color-picker.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import { GenericForm, name, value } from "@climblive/lib/forms";
  import { type Problem } from "@climblive/lib/models";
  import { type Snippet } from "svelte";

  type T = $$Generic<Partial<Problem>>;

  interface Props {
    data: Partial<T>;
    schema: z.ZodType<T, z.ZodTypeDef, T>;
    submit: (value: T) => void;
    children?: Snippet;
  }

  let { data, schema, submit, children }: Props = $props();

  const swatches = [
    "#F44336",
    "#4CAF50",
    "#1790D2",
    "#E410EB",
    "#FFEB3B",
    "#050505",
    "#FF9800",
    "#F628A5",
    "#FAFAFA",
    "#654321",
    "#cccccc",
    "#00FFEF",
  ].join("; ");
</script>

<GenericForm {schema} {submit}>
  <fieldset>
    <wa-input
      size="small"
      {@attach name("number")}
      label="Number"
      type="number"
      required
      {@attach value(data.number)}
    ></wa-input>
    <div class="colors">
      <span>Hold colors</span>
      <div class="pickers">
        <wa-color-picker
          size="small"
          {@attach name("holdColorPrimary")}
          label="Primary hold color"
          required
          {swatches}
          {@attach value(data.holdColorPrimary)}
          no-format-toggle
        ></wa-color-picker>
        <wa-color-picker
          size="small"
          {@attach name("holdColorSecondary")}
          label="Secondary hold color"
          {swatches}
          {@attach value(data.holdColorSecondary)}
          no-format-toggle
        ></wa-color-picker>
      </div>
    </div>
    <wa-input
      size="small"
      {@attach name("description")}
      label="Description"
      type="text"
      {@attach value(data.description)}
    ></wa-input>
    <wa-input
      size="small"
      {@attach name("pointsTop")}
      label="Points for top"
      type="number"
      required
      {@attach value(data.pointsTop)}
    ></wa-input>
    <wa-input
      size="small"
      {@attach name("pointsZone")}
      label="Points for zone"
      type="number"
      required
      {@attach value(data.pointsZone)}
    ></wa-input>
    <wa-input
      size="small"
      {@attach name("flashBonus")}
      label="Flash bonus"
      type="number"
      {@attach value(data.flashBonus)}
    ></wa-input>
    {@render children?.()}
  </fieldset>
</GenericForm>

<style>
  fieldset {
    display: flex;
    flex-direction: column;
    gap: var(--wa-spacing-small);
  }

  .colors {
    span {
      font-size: var(--wa-input-label-font-size-small);
    }

    & .pickers {
      display: flex;
      gap: var(--wa-spacing-x-small);
    }
  }
</style>
