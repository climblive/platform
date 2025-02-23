<script lang="ts">
  import { type ContestTemplate } from "@climblive/lib/models";
  import { serialize } from "@shoelace-style/shoelace";
  import "@shoelace-style/shoelace/dist/components/input/input.js";
  import "@shoelace-style/shoelace/dist/components/option/option.js";
  import "@shoelace-style/shoelace/dist/components/select/select.js";
  import "@shoelace-style/shoelace/dist/components/switch/switch.js";
  import { type Snippet } from "svelte";
  import * as z from "zod";

  const formSchema: z.ZodType<ContestTemplate> = z.object({
    location: z.string().optional(),
    seriesId: z.coerce.number().optional(),
    name: z.string().min(1),
    description: z.string().optional(),
    qualifyingProblems: z.coerce.number().min(0),
    finalists: z.coerce.number().min(0),
    rules: z.string().optional(),
    gracePeriod: z.coerce.number().max(60),
  });

  interface Props {
    data: Partial<ContestTemplate>;
    submit: (patch: ContestTemplate) => void;
    children?: Snippet;
  }

  let { data, submit, children }: Props = $props();

  let form: HTMLFormElement | undefined = $state();

  const handleSubmit = (event: SubmitEvent) => {
    event.preventDefault();

    if (!form) {
      return;
    }

    const data = serialize(form);
    const result = formSchema.safeParse(data);

    if (result.success) {
      submit(result.data);
    } else {
      for (const issue of result.error.issues) {
        setCustomValidity(issue.path, issue.message);
      }
    }

    form?.reportValidity();
  };

  const setCustomValidity = (path: (string | number)[], message: string) => {
    const input = form?.querySelector(`[name="${path}"]`);
    input?.setCustomValidity(message);
  };

  const resetCustomValidation = () => {
    const inputs = form?.querySelectorAll(`[name]`);
    for (const input of inputs) {
      input?.setCustomValidity("");
    }
  };

  const value = (node: HTMLElement, value: string | number | undefined) => {
    $effect(() => {
      node.setAttribute("value", value?.toString() ?? "");
    });
  };

  const name = (node: HTMLElement, value: string | number | undefined) => {
    $effect(() => {
      node.setAttribute("name", value?.toString() ?? "");
    });
  };
</script>

<form bind:this={form} onsubmit={handleSubmit} oninput={resetCustomValidation}>
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
    use:name={"location"}
    label="Location"
    type="text"
    use:value={data.location}
  ></sl-input>
  <sl-input
    size="small"
    use:name={"finalists"}
    label="Finalists"
    help-text="Number of contenders that will proceed to the finals"
    type="number"
    required
    use:value={data.finalists}
    min={0}
    valueAsNumber
  ></sl-input>
  <sl-input
    size="small"
    use:name={"qualifyingProblems"}
    label="Number of qualifying problems"
    help-text="Number of problems that count towards the score"
    type="number"
    required
    use:value={data.qualifyingProblems}
    min={0}
    valueAsNumber
  ></sl-input>
  <sl-input
    size="small"
    use:name={"gracePeriod"}
    label="Grace period"
    help-text="Extra time after the end of the contest during which contenders can enter their last results"
    type="number"
    required
    min={0}
    max={60}
    use:value={data.gracePeriod}
    valueAsNumber
  ></sl-input>
  <sl-textarea
    size="small"
    use:name={"rules"}
    label="Rules"
    use:value={data.rules}
  ></sl-textarea>
  {@render children?.()}
</form>

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: var(--sl-spacing-small);
    padding: var(--sl-spacing-medium);
  }
</style>
