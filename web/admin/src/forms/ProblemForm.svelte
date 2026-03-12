<script lang="ts" module>
  import { z } from "@climblive/lib/utils";

  export const formSchema = z.object({
    number: z.coerce.number().min(0),
    holdColorPrimary: z.string().regex(/^#([0-9a-fA-F]{3}){1,2}$/),
    holdColorSecondary: z.string().optional(),
    description: z.string().optional(),
    zone1Enabled: z.coerce.boolean(),
    zone2Enabled: z.coerce.boolean(),
    pointsZone1: z.coerce
      .number()
      .min(0)
      .max(2 ** 31 - 1)
      .optional(),
    pointsZone2: z.coerce
      .number()
      .min(0)
      .max(2 ** 31 - 1)
      .optional(),
    pointsTop: z.coerce
      .number()
      .min(0)
      .max(2 ** 31 - 1),
    flashBonus: z.coerce
      .number()
      .min(0)
      .max(2 ** 31 - 1)
      .optional(),
  });
</script>

<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/color-picker/color-picker.js";
  import "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import type WaNumberInput from "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import { GenericForm, name } from "@climblive/lib/forms";
  import { type Problem } from "@climblive/lib/models";
  import { type Snippet } from "svelte";

  type T = $$Generic<Partial<Problem>>;

  interface Props {
    data: Partial<T>;
    schema: z.ZodType<T, unknown>;
    submit: (value: T) => void;
    children?: Snippet;
  }

  let { data, schema, submit, children }: Props = $props();

  let zone1Enabled = $state(data.zone1Enabled ?? false);
  let zone2Enabled = $state(data.zone2Enabled ?? false);

  let pointsZone1Input = $state<WaNumberInput>();
  let pointsZone2Input = $state<WaNumberInput>();

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
    "#fff",
  ].join("; ");

  const clearZone1Points = () => {
    if (pointsZone1Input) {
      pointsZone1Input.value = "0";
    }
  };

  const clearZone2Points = () => {
    if (pointsZone2Input) {
      pointsZone2Input.value = "0";
    }
  };

  const addZone1 = () => {
    zone1Enabled = true;
  };

  const removeZone1 = () => {
    zone1Enabled = false;
    clearZone1Points();
    zone2Enabled = false;
    clearZone2Points();
  };

  const addZone2 = () => {
    zone2Enabled = true;
  };

  const removeZone2 = () => {
    zone2Enabled = false;
    clearZone2Points();
  };
</script>

<GenericForm {schema} {submit}>
  <fieldset>
    <wa-number-input
      size="small"
      {@attach name("number")}
      label="Number"
      required
      value={data.number}
      min={0}
    ></wa-number-input>
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

    <div class="card">
      <h3>Top</h3>
      <wa-number-input
        size="small"
        {@attach name("pointsTop")}
        label="Points top"
        hint="Points for reaching the top."
        required
        value={data.pointsTop?.toString() ?? ""}
        min={0}
        max={2 ** 31 - 1}
      >
        <span slot="end">pts</span>
      </wa-number-input>
      <wa-number-input
        size="small"
        {@attach name("flashBonus")}
        label="Flash bonus"
        hint="Bonus points awarded for a flash ascent, added to the total."
        value={data.flashBonus?.toString() ?? ""}
        min={0}
        max={2 ** 31 - 1}
      >
        <span slot="end">pts</span>
      </wa-number-input>
    </div>

    <div class="card">
      <div class="card-header">
        <h3>Z2</h3>
        {#if zone2Enabled}
          <wa-button
            size="small"
            type="button"
            appearance="outlined"
            variant="danger"
            onclick={removeZone2}>Remove</wa-button
          >
        {:else}
          <wa-button
            size="small"
            type="button"
            appearance="outlined"
            onclick={addZone2}
            disabled={!zone1Enabled}>Add</wa-button
          >
        {/if}
      </div>
      <input
        type="hidden"
        name="zone2Enabled"
        value={zone2Enabled ? "on" : ""}
      />
      <wa-number-input
        bind:this={pointsZone2Input}
        size="small"
        {@attach name("pointsZone2")}
        label="Points Z2"
        hint="Points for reaching the second zone."
        value={data.pointsZone2?.toString() ?? ""}
        min={0}
        max={2 ** 31 - 1}
        class={{ hidden: !zone2Enabled }}
      >
        <span slot="end">pts</span>
      </wa-number-input>
    </div>

    <div class="card">
      <div class="card-header">
        <h3>Z1</h3>
        {#if zone1Enabled}
          <wa-button
            size="small"
            type="button"
            appearance="outlined"
            variant="danger"
            onclick={removeZone1}>Remove</wa-button
          >
        {:else}
          <wa-button
            size="small"
            type="button"
            appearance="outlined"
            onclick={addZone1}>Add</wa-button
          >
        {/if}
      </div>
      <input
        type="hidden"
        name="zone1Enabled"
        value={zone1Enabled ? "on" : ""}
      />
      <wa-number-input
        bind:this={pointsZone1Input}
        size="small"
        {@attach name("pointsZone1")}
        label="Points Z1"
        hint="Points for reaching the first zone."
        value={data.pointsZone1?.toString() ?? ""}
        min={0}
        max={2 ** 31 - 1}
        class={{ hidden: !zone1Enabled }}
      >
        <span slot="end">pts</span>
      </wa-number-input>
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

  .card {
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-neutral-border-normal);
    border-radius: var(--wa-border-radius-m);
    padding: var(--wa-space-m);
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);

    & h3 {
      margin: 0;
      font-size: var(--wa-font-size-m);
    }
  }

  .card-header {
    display: flex;
    align-items: center;

    & wa-button {
      margin-inline-start: auto;
    }
  }

  .colors {
    & .pickers {
      display: flex;
      gap: var(--wa-space-s);
    }
  }

  .hidden {
    display: none;
  }
</style>
