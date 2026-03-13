<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import ProblemForm, { formSchema } from "@/forms/ProblemForm.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type { Problem, ProblemPatch } from "@climblive/lib/models";
  import {
    getProblemQuery,
    patchProblemMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  interface Props {
    problemId: number;
  }

  let { problemId }: Props = $props();

  const problemQuery = $derived(getProblemQuery(problemId));
  const patchProblem = $derived(patchProblemMutation(problemId));

  const problem = $derived(problemQuery.data);

  const handleSubmit = async (tmpl: ProblemPatch) => {
    patchProblem.mutate(tmpl, {
      onSuccess: (problem: Problem) =>
        navigate(`/admin/contests/${problem.contestId}#problems`),
      onError: () => toastError("Failed to save problem."),
    });
  };
</script>

{#if problem === undefined}
  <Loader />
{:else}
  <ProblemForm submit={handleSubmit} data={problem} schema={formSchema}>
    <div class="controls">
      <wa-button
        size="small"
        type="button"
        appearance="plain"
        onclick={() =>
          navigate(`/admin/contests/${problem.contestId}#problems`)}
        >Cancel</wa-button
      >
      <wa-button
        size="small"
        type="submit"
        loading={patchProblem.isPending}
        variant="neutral"
        >Save
      </wa-button>
    </div>
  </ProblemForm>
{/if}

<style>
  .controls {
    display: flex;
    gap: var(--wa-space-xs);
    justify-content: right;
  }
</style>
