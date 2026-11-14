<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import CompetitionForm, {
    formSchema,
    nanosecondsInMinute,
  } from "@/forms/CompetitionForm.svelte";
  import type { Competition, CompetitionTemplate } from "@climblive/lib/models";
  import { createCompetitionMutation } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  interface Props {
    organizerId: number;
  }

  let { organizerId }: Props = $props();

  const createCompetition = $derived(createCompetitionMutation(organizerId));

  const handleSubmit = (
    form: Omit<
      CompetitionTemplate,
      | "qualifyingProblems"
      | "finalists"
      | "usePoints"
      | "pooledPoints"
      | "maxAttempts"
      | "pointDeduction"
    >,
  ) => {
    if (createCompetition.isPending) {
      return;
    }

    createCompetition.mutate(
      {
        ...form,
        gracePeriod: form.gracePeriod * nanosecondsInMinute,
        qualifyingProblems: 0,
        finalists: 0,
        usePoints: true,
        pooledPoints: false,
        maxAttempts: 0,
        pointDeduction: 0,
      },
      {
        onSuccess: (competition: Competition) =>
          navigate(`competitions/${competition.id}`),
        onError: () => toastUnexpectedError("Failed to create competition."),
      },
    );
  };
</script>

<CompetitionForm
  submit={handleSubmit}
  data={{
    name: "",
    country: "AQ",
    gracePeriod: 15 * nanosecondsInMinute,
    nameRetentionTime: 43200 * nanosecondsInMinute,
  }}
  schema={formSchema}
>
  <div class="controls">
    <wa-button
      size="s"
      appearance="plain"
      onclick={() => navigate(`./organizers/${organizerId}/competitions`)}
      >Cancel</wa-button
    >
    <wa-button
      size="s"
      type="submit"
      loading={createCompetition.isPending}
      variant="neutral"
      appearance="accent"
      >Create
    </wa-button>
  </div>
</CompetitionForm>

<style>
  .controls {
    margin-inline-start: auto;
  }
</style>
