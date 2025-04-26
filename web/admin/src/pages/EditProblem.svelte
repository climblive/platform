<script lang="ts">
  import ProblemForm from "@/forms/ProblemForm.svelte";
  import type { ProblemPatch } from "@climblive/lib/models";
  import {
    getProblemQuery,
    patchProblemMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import * as z from "zod";

  interface Props {
    problemId: number;
  }

  let { problemId }: Props = $props();

  const formSchema: z.ZodType<ProblemPatch> = z.object({
    number: z.coerce.number(),
    holdColorPrimary: z.string().regex(/^#([0-9a-fA-F]{3}){1,2}$/),
    holdColorSecondary: z.string().optional(),
    description: z.string().optional(),
    pointsTop: z.coerce.number(),
    pointsZone: z.coerce.number(),
    flashBonus: z.coerce.number().optional(),
  });

  const problemQuery = getProblemQuery(problemId);
  const patchProblem = patchProblemMutation(problemId);

  const problem = $derived($problemQuery.data);

  const handleSubmit = async (tmpl: ProblemPatch) => {
    $patchProblem.mutate(tmpl, {
      onError: () => toastError("Failed to save problem."),
    });
  };
</script>

<ProblemForm
  submit={handleSubmit}
  data={{
    ...problem,
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
      loading={$patchProblem.isPending}
      variant="primary"
      >Save
    </sl-button>
  </div>
</ProblemForm>

<style>
  .controls {
    display: flex;
    justify-content: start;
  }
</style>
