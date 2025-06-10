<script lang="ts">
  import ProblemForm, { formSchema } from "@/forms/ProblemForm.svelte";
  import type { Problem, ProblemPatch } from "@climblive/lib/models";
  import {
    getProblemQuery,
    patchProblemMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import { navigate } from "svelte-routing";

  interface Props {
    problemId: number;
  }

  let { problemId }: Props = $props();

  const problemQuery = $derived(getProblemQuery(problemId));
  const patchProblem = $derived(patchProblemMutation(problemId));

  const problem = $derived($problemQuery.data);

  const handleSubmit = async (tmpl: ProblemPatch) => {
    $patchProblem.mutate(tmpl, {
      onSuccess: (problem: Problem) =>
        navigate(`/admin/contests/${problem.contestId}#problems`),
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
