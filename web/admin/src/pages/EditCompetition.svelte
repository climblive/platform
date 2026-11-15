<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import CompetitionForm, {
    editFormSchema,
    nanosecondsInMinute,
  } from "@/forms/CompetitionForm.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type { Competition, CompetitionPatch } from "@climblive/lib/models";
  import {
    getCompetitionQuery,
    patchCompetitionMutation,
  } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  interface Props {
    competitionId: number;
  }

  let { competitionId }: Props = $props();

  const competitionQuery = $derived(getCompetitionQuery(competitionId));
  const patchCompetition = $derived(patchCompetitionMutation(competitionId));

  const competition = $derived(competitionQuery.data);

  const handleSubmit = async (tmpl: CompetitionPatch) => {
    patchCompetition.mutate(
      {
        ...tmpl,
        gracePeriod:
          tmpl.gracePeriod !== undefined
            ? tmpl.gracePeriod * nanosecondsInMinute
            : undefined,
      },
      {
        onSuccess: (competition: Competition) =>
          navigate(`/admin/competitions/${competition.id}`),
        onError: () => toastUnexpectedError("Failed to save competition."),
      },
    );
  };
</script>

{#if competition === undefined}
  <Loader />
{:else}
  <CompetitionForm
    submit={handleSubmit}
    data={competition}
    schema={editFormSchema}
    disableNameRetentionTime
  >
    <div class="controls">
      <wa-button
        size="s"
        type="button"
        appearance="plain"
        onclick={history.back()}>Cancel</wa-button
      >
      <wa-button
        size="s"
        type="submit"
        loading={patchCompetition.isPending}
        variant="neutral"
        >Save
      </wa-button>
    </div>
  </CompetitionForm>
{/if}

<style>
  .controls {
    display: flex;
    justify-content: end;
    gap: var(--wa-space-xs);
  }
</style>
