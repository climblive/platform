<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import ProblemForm, { formSchema } from "@/forms/ProblemForm.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type { ProblemTemplate } from "@climblive/lib/models";
  import {
    createProblemMutation,
    getProblemsQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const problemsQuery = $derived(getProblemsQuery(contestId));

  let highestProblemNumber = $derived.by(() => {
    if ($problemsQuery.data === undefined) {
      return undefined;
    } else if ($problemsQuery.data.length > 0) {
      return Math.max(
        ...($problemsQuery.data?.map(({ number }) => number) ?? []),
      );
    } else {
      return 0;
    }
  });

  const createProblem = $derived(createProblemMutation(contestId));

  const handleSubmit = async (tmpl: Omit<ProblemTemplate, "pointsZone">) => {
    $createProblem.mutate(
      { ...tmpl, pointsZone: 0 },
      {
        onSuccess: () => navigate(`/admin/contests/${contestId}#contest`),
        onError: () => toastError("Failed to create problem."),
      },
    );
  };
</script>

{#if highestProblemNumber === undefined}
  <Loader />
{:else}
  <ProblemForm
    submit={handleSubmit}
    data={{
      number: highestProblemNumber + 1,
      holdColorPrimary: "#000000",
      pointsTop: 100,
      flashBonus: 0,
    }}
    schema={formSchema}
  >
    <div class="controls">
      <wa-button
        size="small"
        type="button"
        appearance="plain"
        onclick={() => navigate(`/admin/contests/${contestId}#problems`)}
        >Cancel</wa-button
      >
      <wa-button
        size="small"
        type="submit"
        loading={$createProblem.isPending}
        variant="neutral"
        >Create
      </wa-button>
    </div>
  </ProblemForm>
{/if}
