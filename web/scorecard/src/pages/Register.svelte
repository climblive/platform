<script lang="ts">
  import RegistrationForm from "@/forms/RegistrationForm.svelte";
  import type { ScorecardSession } from "@/types";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import {
    ContestStateProvider,
    SplashScreen,
  } from "@climblive/lib/components";
  import type { ContenderPatch } from "@climblive/lib/models";
  import {
    getContenderQuery,
    getContestQuery,
    patchContenderMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { add, formatDistance } from "date-fns";
  import { getContext } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Readable } from "svelte/store";

  const nanosecondsInMinute = 60 * 1_000_000_000;

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  const contenderQuery = $derived(getContenderQuery($session.contenderId));
  const contestQuery = $derived(getContestQuery($session.contestId));
  const patchContender = $derived(patchContenderMutation($session.contenderId));

  const contender = $derived(contenderQuery.data);
  const contest = $derived(contestQuery.data);

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
        onError: () => toastError("Registration was not successful."),
      },
    );
  };

  const retentionDuration = $derived.by(() => {
    const base = new Date(0);
    return formatDistance(
      add(base, {
        minutes: (contest?.nameRetentionTime ?? 0) / nanosecondsInMinute,
      }),
      base,
    );
  });

  let showSplash = $state(true);
</script>

{#if showSplash || !contender || !contest}
  <SplashScreen onComplete={() => (showSplash = false)} />
{:else}
  <h1>{contest.name}</h1>
  <ContestStateProvider contestId={contest.id}>
    {#snippet children({ contestState })}
      <RegistrationForm
        submit={handleSubmit}
        data={{
          name: contender.name,
          compClassId: contender.compClassId,
          withdrawnFromFinals: contender.withdrawnFromFinals,
        }}
        callout={registerCallout}
        {contestState}
      >
        <wa-button
          size="small"
          type="submit"
          loading={patchContender.isPending}
          variant="neutral"
          appearance="accent"
          >Register
        </wa-button>
      </RegistrationForm>

      {#snippet registerCallout()}
        <wa-callout variant="neutral" size="small">
          <wa-icon slot="icon" name="circle-info"></wa-icon>
          Your name will be stored for {retentionDuration} after the contest ends,
          after which it will be removed and your results anonymized.
        </wa-callout>
      {/snippet}
    {/snippet}
  </ContestStateProvider>
{/if}

<style>
  h1 {
    font-size: var(--wa-font-size-l);
    padding: var(--wa-space-m);
    padding-block-end: 0;
  }
</style>
