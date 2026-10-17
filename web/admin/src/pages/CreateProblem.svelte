<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import ProblemForm, { formSchema } from "@/forms/ProblemForm.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type { Problem, ProblemTemplate } from "@climblive/lib/models";
  import {
    createProblemMutation,
    getContestQuery,
    getProblemsQuery,
  } from "@climblive/lib/queries";
  import {
    compareLabels,
    labelRegex,
    toastUnexpectedError,
  } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const contestQuery = $derived(getContestQuery(contestId));
  const contest = $derived(contestQuery.data);

  const problemsQuery = $derived(getProblemsQuery(contestId));

  let highestProblem = $derived.by<Problem | undefined | null>(() => {
    if (problemsQuery.data === undefined) {
      return undefined;
    } else if (problemsQuery.data.length > 0) {
      const sortedProblems = problemsQuery.data.toSorted((p1, p2) =>
        compareLabels(p1.label, p2.label),
      );

      return sortedProblems.slice(-1)[0];
    } else {
      return null;
    }
  });

  const createProblem = $derived(createProblemMutation(contestId));

  const handleSubmit = async (tmpl: Omit<ProblemTemplate, "pointsZone">) => {
    createProblem.mutate(
      { ...tmpl },
      {
        onSuccess: () => navigate(`/admin/contests/${contestId}#problems`),
        onError: () => toastUnexpectedError("Failed to create problem."),
      },
    );
  };

  const template = $derived.by<Partial<Problem>>(() => {
    if (highestProblem) {
      let label = "0";

      const matches = highestProblem.label.match(labelRegex);
      if (matches !== null) {
        label = matches[1] + (Number(matches[2]) + 1) + matches[3];
      }

      return {
        ...highestProblem,
        label,
        holdColorPrimary: "#000000",
        holdColorSecondary: undefined,
      };
    } else {
      return {
        label: "1",
        holdColorPrimary: "#000000",
        pointsTop: contest?.usePoints ? 100 : 0,
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
    hidePoints={contest?.usePoints === false}
  >
    <div class="controls">
      <wa-button
        size="s"
        type="button"
        appearance="plain"
        onclick={() => navigate(`/admin/contests/${contestId}#problems`)}
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
