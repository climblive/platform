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
  import "@awesome.me/webawesome/dist/components/divider/divider.js";
  import "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import type WaNumberInput from "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import type WaSwitch from "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { HoldColorPicker } from "@climblive/lib/components";
  import { checked, GenericForm, name } from "@climblive/lib/forms";
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

  let zone1Enabled = $derived(data.zone1Enabled);
  let zone2Enabled = $derived(data.zone2Enabled);

  let pointsZone1Input = $state<WaNumberInput>();
  let pointsZone2Input = $state<WaNumberInput>();

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

  const handleZone1Toggle = (event: InputEvent) => {
    const target = event.target as WaSwitch;
    zone1Enabled = target.checked;

    if (!zone1Enabled) {
      clearZone1Points();

      zone2Enabled = false;
      clearZone2Points();
    }
  };

  const handleZone2Toggle = (event: InputEvent) => {
    const target = event.target as WaSwitch;
    zone2Enabled = target.checked;

    if (!zone2Enabled) {
      clearZone2Points();
    }
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
      <HoldColorPicker
        name="holdColorPrimary"
        label="Primary hold color"
        required
        value={data.holdColorPrimary}
        placement="right"
      />
      <HoldColorPicker
        name="holdColorSecondary"
        label="Secondary hold color"
        allowClear={true}
        value={data.holdColorSecondary}
        placement="bottom"
      />
    </div>
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

    <wa-divider></wa-divider>

    <wa-switch
      size="small"
      {@attach name("zone1Enabled")}
      hint="Add a zone."
      onchange={handleZone1Toggle}
      {@attach checked(zone1Enabled)}>Enable zone Z1</wa-switch
    >
    <wa-number-input
      bind:this={pointsZone1Input}
      size="small"
      {@attach name("pointsZone1")}
      label="Points Z1"
      hint="Points for reaching the first zone."
      value={data.pointsZone1?.toString() ?? ""}
      min={0}
      max={2 ** 31 - 1}
      class={{
        hidden: !zone1Enabled,
      }}
    >
      <span slot="end">pts</span>
    </wa-number-input>
    {#if zone1Enabled}
      <wa-switch
        size="small"
        {@attach name("zone2Enabled")}
        hint="Add a second zone."
        onchange={handleZone2Toggle}
        {@attach checked(zone2Enabled)}>Enable zone Z2</wa-switch
      >
    {/if}
    <wa-number-input
      bind:this={pointsZone2Input}
      size="small"
      {@attach name("pointsZone2")}
      label="Points Z2"
      hint="Points for reaching the second zone."
      value={data.pointsZone2?.toString() ?? ""}
      min={0}
      max={2 ** 31 - 1}
      class={{
        hidden: !zone2Enabled,
      }}
    >
      <span slot="end">pts</span>
    </wa-number-input>

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
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--wa-space-s);
  }

  .hidden {
    display: none;
  }
</style>
