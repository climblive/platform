<script lang="ts" module>
  import * as z from "zod";

  const formSchema = z.object({
    location: z.string().optional(),
    seriesId: z.coerce.number().optional(),
    name: z.string().min(1),
    description: z.string().optional(),
    qualifyingProblems: z.coerce.number().min(0),
    finalists: z.coerce.number().min(0),
    rules: z.string().optional(),
    gracePeriod: z.coerce.number().min(0).max(60),
  });
</script>

<script lang="ts">
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import { GenericForm, name } from "@climblive/lib/forms";
  import { type ContestTemplate } from "@climblive/lib/models";
  import { type Snippet } from "svelte";

  interface Props {
    data: Partial<ContestTemplate>;
    submit: (value: ContestTemplate) => void;
    children?: Snippet;
  }

  let { data, submit, children }: Props = $props();
</script>

<GenericForm schema={formSchema} {submit}>
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
    ></wa-input>
    <wa-input
      size="small"
      {@attach name("finalists")}
      label="Finalists"
      hint="Number of contenders that will proceed to the finals"
      type="number"
      required
      value={data.finalists}
      min={0}
    ></wa-input>
    <wa-input
      size="small"
      {@attach name("qualifyingProblems")}
      label="Number of qualifying problems"
      hint="Number of problems that count towards the score"
      type="number"
      required
      value={data.qualifyingProblems}
      min={0}
    ></wa-input>
    <wa-input
      size="small"
      {@attach name("gracePeriod")}
      label="Grace period"
      hint="Extra time after the end of the contest during which contenders can enter their last results"
      type="number"
      required
      min={0}
      max={60}
      value={data.gracePeriod}
    ></wa-input>
    <wa-textarea
      size="small"
      {@attach name("rules")}
      label="Rules"
      value={data.rules}
    ></wa-textarea>
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
