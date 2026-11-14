<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import ProblemForm, { formSchema } from "@/forms/ProblemForm.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type { Problem, ProblemTemplate } from "@climblive/lib/models";
  import {
    createProblemMutation,
    getCompetitionQuery,
    getProblemsQuery,
  } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  interface Props {
    competitionId: number;
  }

  let { competitionId }: Props = $props();

  const competitionQuery = $derived(getCompetitionQuery(competitionId));
  const competition = $derived(competitionQuery.data);

  const problemsQuery = $derived(getProblemsQuery(competitionId));

  let highestProblem = $derived.by<Problem | undefined | null>(() => {
    if (problemsQuery.data === undefined) {
      return undefined;
    } else if (problemsQuery.data.length > 0) {
      const sortedProblems = problemsQuery.data.toSorted(
        (p1, p2) => p1.number - p2.number,
      );

      return sortedProblems.slice(-1)[0];
    } else {
      return null;
    }
  });

  const createProblem = $derived(createProblemMutation(competitionId));

  const handleSubmit = async (tmpl: Omit<ProblemTemplate, "pointsZone">) => {
    createProblem.mutate(
      { ...tmpl },
      {
        onSuccess: () =>
          navigate(`/admin/competitions/${competitionId}#problems`),
        onError: () => toastUnexpectedError("Failed to create problem."),
      },
    );
  };

  const template = $derived.by<Partial<Problem>>(() => {
    if (highestProblem) {
      return {
        ...highestProblem,
        number: highestProblem.number + 1,
        holdColorPrimary: "#000000",
        holdColorSecondary: undefined,
      };
    } else {
      return {
        number: 1,
        holdColorPrimary: "#000000",
        pointsTop: competition?.usePoints ? 100 : 0,
        flashBonus: 0,
      };
    }
  });
</script>

{#if highestProblem === undefined}
  <Loader />
{:else}
  <ProblemForm
    submit={handleSubmit}
    data={template}
    schema={formSchema}
    hidePoints={competition?.usePoints === false}
  >
    <div class="controls">
      <wa-button
        size="s"
        type="button"
        appearance="plain"
        onclick={() =>
          navigate(`/admin/competitions/${competitionId}#problems`)}
        >Cancel</wa-button
      >
      <wa-button
        size="s"
        type="submit"
        loading={createProblem.isPending}
        variant="neutral"
        >Create
      </wa-button>
    </div>
  </ProblemForm>
{/if}

<style>
  .controls {
    display: flex;
    gap: var(--wa-space-xs);
    justify-content: end;
  }
</style>
