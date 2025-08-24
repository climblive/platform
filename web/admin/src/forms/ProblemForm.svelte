<script lang="ts" module>
  import * as z from "zod";

  export const formSchema = z.object({
    number: z.coerce.number(),
    holdColorPrimary: z.string().regex(/^#([0-9a-fA-F]{3}){1,2}$/),
    holdColorSecondary: z.string().optional(),
    pointsTop: z.coerce.number(),
    flashBonus: z.coerce.number().optional(),
  });
</script>

<script lang="ts">
  import "@awesome.me/webawesome/dist/components/color-picker/color-picker.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import { GenericForm, name } from "@climblive/lib/forms";
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
      value={data.number}
    ></wa-input>
    <div class="colors">
      <div class="pickers">
        <wa-color-picker
          size="small"
          {@attach name("holdColorPrimary")}
          label="Primary hold color"
          required
          {swatches}
          value={data.holdColorPrimary}
          without-format-toggle
        ></wa-color-picker>
        <wa-color-picker
          size="small"
          {@attach name("holdColorSecondary")}
          label="Secondary hold color"
          {swatches}
          value={data.holdColorSecondary}
          without-format-toggle
        ></wa-color-picker>
      </div>
    </div>
    <wa-input
      size="small"
      {@attach name("pointsTop")}
      label="Points for top"
      type="number"
      required
      value={data.pointsTop?.toString() ?? ""}
    ></wa-input>
    <wa-input
      size="small"
      {@attach name("flashBonus")}
      label="Flash bonus"
      type="number"
      value={data.flashBonus?.toString() ?? ""}
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

  .colors {
    & .pickers {
      display: flex;
      gap: var(--wa-space-s);
    }
  }
</style>
