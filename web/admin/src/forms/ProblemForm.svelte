<script lang="ts">
  import { GenericForm, name, value } from "@climblive/lib/forms";
  import { type Problem } from "@climblive/lib/models";
  import "@shoelace-style/shoelace/dist/components/color-picker/color-picker.js";
  import "@shoelace-style/shoelace/dist/components/input/input.js";
  import { type Snippet } from "svelte";
  import * as z from "zod";

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
    <sl-input
      size="small"
      use:name={"number"}
      label="Number"
      type="number"
      required
      use:value={data.number}
    ></sl-input>
    <div class="colors">
      <span>Hold colors</span>
      <div class="pickers">
        <sl-color-picker
          size="small"
          use:name={"holdColorPrimary"}
          label="Primary hold color"
          required
          {swatches}
          use:value={data.holdColorPrimary}
          no-format-toggle
        ></sl-color-picker>
        <sl-color-picker
          size="small"
          use:name={"holdColorSecondary"}
          label="Secondary hold color"
          {swatches}
          use:value={data.holdColorSecondary}
          no-format-toggle
        ></sl-color-picker>
      </div>
    </div>
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

  .colors {
    span {
      font-size: var(--sl-input-label-font-size-small);
    }

    & .pickers {
      display: flex;
      gap: var(--sl-spacing-x-small);
    }
  }
</style>
