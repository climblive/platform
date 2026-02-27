<script lang="ts" module>
  import { z } from "@climblive/lib/utils";

  export const formSchema = z.object({
    location: z.string().optional(),
    country: z.string(),
    seriesId: z.coerce.number().optional(),
    name: z.string().min(1),
    description: z.string().optional(),
    info: z.string().optional(),
    gracePeriod: z.coerce.number().min(0).max(60),
  });

  export const minuteInNanoseconds = 60 * 1_000_000_000;
</script>

<script lang="ts">
  import InfoInput from "@/components/InfoInput.svelte";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import "@awesome.me/webawesome/dist/components/textarea/textarea.js";
  import { GenericForm, name, value } from "@climblive/lib/forms";
  import type { Contest } from "@climblive/lib/models";
  import { countries, getFlag } from "@climblive/lib/utils";
  import { type Snippet } from "svelte";

  type T = $$Generic<Partial<Contest>>;

  interface Props {
    data: Partial<T>;
    schema: z.ZodType<T, unknown>;
    submit: (value: T) => void;
    children?: Snippet;
  }

  let { data, schema, submit, children }: Props = $props();

  let selectedCountry = $derived(data.country || "AQ");

  const handleCountryChange = (event: Event) => {
    const target = event.target as WaSelect;

    if (typeof target.value === "string") {
      selectedCountry = target.value;
    }
  };

  let showGeneralInfo = $derived(data.info !== undefined);
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
    ></wa-input>
    <wa-input
      size="small"
      {@attach name("description")}
      label="Description"
      type="text"
      value={data.description}
    ></wa-input>
    <div class="location">
      <wa-input
        size="small"
        {@attach name("location")}
        label="Location"
        type="text"
        value={data.location}
        hint="Usually the name of the climbing gym."
      ></wa-input>
      <wa-select
        size="small"
        {@attach name("country")}
        {@attach value(selectedCountry)}
        label="Country"
        onchange={handleCountryChange}
      >
        <span slot="start">{getFlag(selectedCountry)}</span>
        {#each countries as country (country.code)}
          <wa-option value={country.code} label={country.name}>
            <span slot="start">{getFlag(country.code)}</span>
            {country.name}
          </wa-option>
        {/each}
      </wa-select>
    </div>
    <wa-number-input
      size="small"
      {@attach name("gracePeriod")}
      label="Grace period"
      hint="Extra time after the end of the contest during which contenders can enter their last results."
      required
      min={0}
      max={60}
      value={Math.floor((data.gracePeriod ?? 0) / minuteInNanoseconds)}
    >
      <span slot="end">minutes</span>
    </wa-number-input>
    {#if showGeneralInfo}
      <InfoInput info={data.info} />
    {:else}
      <wa-button
        size="small"
        appearance="outlined"
        onclick={() => (showGeneralInfo = true)}
      >
        <wa-icon slot="start" name="plus"></wa-icon>
        Add general info
      </wa-button>
    {/if}
    {@render children?.()}
  </fieldset>
</GenericForm>

<style>
  fieldset {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }

  .location {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--wa-space-s);
  }

  @media screen and (max-width: 768px) {
    .location {
      grid-template-columns: 1fr;
    }
  }

  wa-button {
    align-self: start;
  }
</style>
