<script lang="ts">
  import RegistrationForm from "@/forms/RegistrationForm.svelte";
  import type { ScorecardSession } from "@/types";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import {
    CompetitionStateProvider,
    SplashScreen,
  } from "@climblive/lib/components";
  import type { ContenderPatch } from "@climblive/lib/models";
  import {
    getCategoriesQuery,
    getContenderQuery,
    getCompetitionQuery,
    patchContenderMutation,
  } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import { add, formatDistance, isAfter } from "date-fns";
  import { getContext } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Readable } from "svelte/store";

  const nanosecondsInMinute = 60 * 1_000_000_000;

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  const contenderQuery = $derived(getContenderQuery($session.contenderId));
  const competitionQuery = $derived(
    getCompetitionQuery($session.competitionId),
  );
  const categoriesQuery = $derived(getCategoriesQuery($session.competitionId));
  const patchContender = $derived(patchContenderMutation($session.contenderId));

  const contender = $derived(contenderQuery.data);
  const competition = $derived(competitionQuery.data);

  const tooLate = $derived(
    categoriesQuery.data?.every(({ timeEnd }) => isAfter(new Date(), timeEnd)),
  );

  const gotoScorecard = () => {
    navigate(`/${contender?.registrationCode}`, { replace: true });
  };

  const handleSubmit = (form: ContenderPatch) => {
    if (!contender || patchContender.isPending) {
      return;
    }

    patchContender.mutate(
      {
        ...form,
      },
      {
        onSuccess: gotoScorecard,
        onError: () => toastUnexpectedError("Registration was not successful."),
      },
    );
  };

  const retentionDuration = $derived.by(() => {
    const base = new Date(0);
    return formatDistance(
      add(base, {
        minutes: (competition?.nameRetentionTime ?? 0) / nanosecondsInMinute,
      }),
      base,
    );
  });

  let showSplash = $state(true);
</script>

{#if showSplash || !contender || !competition || tooLate === undefined}
  <SplashScreen onComplete={() => (showSplash = false)} />
{:else}
  <h1>{competition.name}</h1>
  <CompetitionStateProvider competitionId={competition.id}>
    {#snippet children({ competitionState })}
      <RegistrationForm
        submit={handleSubmit}
        data={{
          name: contender.name,
          categoryId: contender.categoryId,
          withdrawnFromFinals: contender.withdrawnFromFinals,
        }}
        callout={registerCallout}
        {competitionState}
      >
        {#if tooLate}
          <wa-callout variant="warning">
            <wa-icon slot="icon" name="clock"></wa-icon>
            <strong>Registration is no longer possible</strong><br />
            All categories have ended.
          </wa-callout>
        {/if}
        <wa-button
          size="s"
          type="submit"
          loading={patchContender.isPending}
          variant="neutral"
          appearance="accent"
          disabled={tooLate}
          >Register
        </wa-button>
      </RegistrationForm>

      {#snippet registerCallout()}
        <wa-callout variant="neutral" size="s">
          <wa-icon slot="icon" name="circle-info"></wa-icon>
          Your name will be stored for {retentionDuration} after the competition ends,
          after which it will be removed and your results anonymized.
        </wa-callout>
      {/snippet}
    {/snippet}
  </CompetitionStateProvider>
{/if}

<style>
  h1 {
    font-size: var(--wa-font-size-l);
    padding: var(--wa-space-m);
    padding-block-end: 0;
  }
</style>
