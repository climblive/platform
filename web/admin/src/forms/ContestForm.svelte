<script lang="ts" module>
  import * as z from "zod/v4";

  export const formSchema = z.object({
    location: z.string().optional(),
    country: z.string().optional(),
    seriesId: z.coerce.number().optional(),
    name: z.string().min(1),
    description: z.string().optional(),
    info: z.string().optional(),
    gracePeriod: z.coerce.number().min(0).max(60),
  });

  export const minuteInNanoseconds = 60 * 1_000_000_000;

  export const countries = [
    { code: "SE", name: "Sweden" },
    { code: "NO", name: "Norway" },
    { code: "DK", name: "Denmark" },
    { code: "FI", name: "Finland" },
    { code: "IS", name: "Iceland" },
    { code: "DE", name: "Germany" },
    { code: "FR", name: "France" },
    { code: "GB", name: "United Kingdom" },
    { code: "ES", name: "Spain" },
    { code: "IT", name: "Italy" },
    { code: "NL", name: "Netherlands" },
    { code: "BE", name: "Belgium" },
    { code: "AT", name: "Austria" },
    { code: "CH", name: "Switzerland" },
    { code: "PL", name: "Poland" },
    { code: "CZ", name: "Czech Republic" },
    { code: "US", name: "United States" },
    { code: "CA", name: "Canada" },
    { code: "AU", name: "Australia" },
    { code: "NZ", name: "New Zealand" },
    { code: "JP", name: "Japan" },
  ];

  function getFlag(countryCode: string): string {
    if (!countryCode || countryCode.length !== 2) return "";
    const codePoints = countryCode
      .toUpperCase()
      .split("")
      .map((char) => 127397 + char.charCodeAt(0));
    return String.fromCodePoint(...codePoints);
  }
</script>

<script lang="ts">
  import InfoInput from "@/components/InfoInput.svelte";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import "@awesome.me/webawesome/dist/components/textarea/textarea.js";
  import { GenericForm, name } from "@climblive/lib/forms";
  import type { Contest } from "@climblive/lib/models";
  import { type Snippet } from "svelte";

  type T = $$Generic<Partial<Contest>>;

  interface Props {
    data: Partial<T>;
    schema: z.ZodType<T, unknown>;
    submit: (value: T) => void;
    children?: Snippet;
  }

  let { data, schema, submit, children }: Props = $props();
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
      label="Country"
      value={data.country}
      clearable
    >
      {#each countries as country (country.code)}
        <wa-option value={country.code}>
          {getFlag(country.code)}
          {country.name}
        </wa-option>
      {/each}
    </wa-select>
    <wa-input
      size="small"
      {@attach name("gracePeriod")}
      label="Grace period (minutes)"
      hint="Extra time after the end of the contest during which contenders can enter their last results."
      type="number"
      required
      min={0}
      max={60}
      value={Math.floor((data.gracePeriod ?? 0) / minuteInNanoseconds)}
    >
    </wa-input>
    <InfoInput info={data.info} />
    {@render children?.()}
  </fieldset>
</GenericForm>

<style>
  fieldset {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }
</style>
