<script lang="ts">
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import RegistrationForm from "@/forms/RegistrationForm.svelte";
  import type { ScorecardSession } from "@/types";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { CompetitionStateProvider } from "@climblive/lib/components";
  import type { CompetitorPatch } from "@climblive/lib/models";
  import {
    getCompetitorQuery,
    getCompetitionQuery,
    patchCompetitorMutation,
    scrubCompetitorMutation,
  } from "@climblive/lib/queries";
  import { SyncedTime, toastUnexpectedError } from "@climblive/lib/utils";
  import { formatDistance, isBefore } from "date-fns";
  import { getContext, onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Readable } from "svelte/store";
  import Loading from "./Loading.svelte";

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  const competitorQuery = $derived(getCompetitorQuery($session.competitorId));
  const competitionQuery = $derived(
    getCompetitionQuery($session.competitionId),
  );
  const patchCompetitor = $derived(
    patchCompetitorMutation($session.competitorId),
  );
  const scrubCompetitor = $derived(
    scrubCompetitorMutation($session.competitorId),
  );
  const time = new SyncedTime(60_000);

  onMount(() => {
    time.start();

    return () => time.stop();
  });

  let competitor = $derived(competitorQuery.data);
  let competition = $derived(competitionQuery.data);

  const retentionDuration = $derived.by(() => {
    if (!competitor?.scrubBefore) {
      return undefined;
    }

    const now = time.current;

    if (isBefore(competitor.scrubBefore, now)) {
      return undefined;
    }

    return formatDistance(competitor.scrubBefore, now);
  });

  let scrubDialog: WaDialog | undefined = $state();

  const gotoScorecard = () => {
    navigate(`/${competitor?.registrationCode}`);
  };

  const handleSubmit = (form: CompetitorPatch) => {
    if (!competitor || patchCompetitor.isPending) {
      return;
    }

    patchCompetitor.mutate(
      {
        ...form,
      },
      {
        onSuccess: gotoScorecard,
        onError: () =>
          toastUnexpectedError("Failed to save registration data."),
      },
    );
  };

  const handleScrub = () => {
    scrubCompetitor.mutate(undefined, {
      onSuccess: gotoScorecard,
      onError: () => toastUnexpectedError("Failed to remove your name."),
    });
  };
</script>

{#if !competitor || !competition}
  <Loading />
{:else}
  <CompetitionStateProvider
    competitionId={competition.id}
    categoryId={competitor.categoryId}
  >
    {#snippet children({ competitionState })}
      {@const disabled =
        competitor.scrubbedAt !== undefined || competitionState === "ENDED"}

      <RegistrationForm
        submit={handleSubmit}
        data={{
          name: competitor.name,
          categoryId: competitor.categoryId,
          withdrawnFromFinals: competitor.withdrawnFromFinals,
        }}
        callout={profileCallout}
        {competitionState}
        {disabled}
      >
        <div class="controls">
          {#if !!competitor.name}
            <wa-button
              class="scrub-button"
              size="s"
              type="button"
              variant="danger"
              appearance="outlined"
              onclick={() => {
                if (scrubDialog) {
                  scrubDialog.open = true;
                }
              }}
            >
              <wa-icon slot="start" name="user-secret"></wa-icon>
              Remove my name
            </wa-button>
          {/if}
          <wa-button
            size="s"
            type="button"
            onclick={gotoScorecard}
            appearance="plain"
            >Cancel
          </wa-button>
          <wa-button
            size="s"
            type="submit"
            loading={patchCompetitor.isPending}
            variant="neutral"
            appearance="accent"
            {disabled}
            >Save
          </wa-button>
        </div>
      </RegistrationForm>

      {#snippet profileCallout()}
        {#if !competitor?.scrubbedAt}
          <wa-callout variant="neutral" size="s">
            <wa-icon slot="icon" name="circle-info"></wa-icon>
            {#if retentionDuration}
              Your name will be kept stored for {retentionDuration} from now, after
              which it will be removed and your results anonymized.
            {:else}
              Your name will be removed and your results anonymized shortly.
            {/if}
          </wa-callout>
        {/if}
      {/snippet}
    {/snippet}
  </CompetitionStateProvider>

  <wa-dialog bind:this={scrubDialog} label="Remove your name">
    Your name will be <strong>permanently</strong> removed and your results will
    be anonymized. This action cannot be undone.
    <br /><br />
    Be aware that without a name, you will lose your chance at finals and you cannot
    take part in any prize raffles.
    <wa-button
      size="s"
      slot="footer"
      appearance="plain"
      onclick={() => {
        if (scrubDialog) {
          scrubDialog.open = false;
        }
      }}>Cancel</wa-button
    >
    <wa-button
      size="s"
      slot="footer"
      variant="danger"
      loading={scrubCompetitor.isPending}
      onclick={handleScrub}
    >
      <wa-icon slot="start" name="user-secret"></wa-icon>
      Remove my name anyway
    </wa-button>
  </wa-dialog>
{/if}

<style>
  .controls {
    display: flex;
    justify-content: end;
    gap: var(--wa-space-xs);

    & > .scrub-button {
      margin-right: auto;
    }
  }
</style>
