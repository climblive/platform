<script lang="ts">
  import ProblemForm from "@/forms/ProblemForm.svelte";
  import type { ProblemTemplate } from "@climblive/lib/models";
  import { createProblemMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import * as z from "zod";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const formSchema: z.ZodType<ProblemTemplate> = z.object({
    number: z.coerce.number(),
    holdColorPrimary: z.string().regex(/^#([0-9a-fA-F]{3}){1,2}$/),
    holdColorSecondary: z.string().optional(),
    description: z.string().optional(),
    pointsTop: z.coerce.number(),
    pointsZone: z.coerce.number(),
    flashBonus: z.coerce.number().optional(),
  });

  const createProblem = createProblemMutation(contestId);

  const handleSubmit = async (tmpl: ProblemTemplate) => {
    $createProblem.mutate(tmpl, {
      onError: () => toastError("Failed to create problem."),
    });
  };
</script>

<ProblemForm
  submit={handleSubmit}
  data={{
    number: 1,
    holdColorPrimary: "#000000",
    pointsTop: 100,
    pointsZone: 0,
    flashBonus: 0,
  }}
  schema={formSchema}
>
  <div class="controls">
    <sl-button
      size="small"
      type="button"
      variant="text"
      onclick={history.back()}>Cancel</sl-button
    >
    <sl-button
      size="small"
      type="submit"
      loading={$createProblem.isPending}
      variant="primary"
      >Create
    </sl-button>
  </div>
</ProblemForm>
