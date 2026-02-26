<script lang="ts">
  import RegistrationForm from "@/forms/RegistrationForm.svelte";
  import type { ScorecardSession } from "@/types";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { ContestStateProvider } from "@climblive/lib/components";
  import type { ContenderPatch } from "@climblive/lib/models";
  import {
    getContenderQuery,
    getContestQuery,
    patchContenderMutation,
    scrubContenderMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { add } from "date-fns";
  import { getContext } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Readable } from "svelte/store";
  import Loading from "./Loading.svelte";

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  const contenderQuery = $derived(getContenderQuery($session.contenderId));
  const contestQuery = $derived(getContestQuery($session.contestId));
  const patchContender = $derived(patchContenderMutation($session.contenderId));
  const scrubContender = $derived(scrubContenderMutation($session.contenderId));

  let contender = $derived(contenderQuery.data);
  let contest = $derived(contestQuery.data);

  let scrubDialog: WaDialog | undefined = $state();

  const gotoScorecard = () => {
    navigate(`/${contender?.registrationCode}`);
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
        onError: () => toastError("Failed to save registration data."),
      },
    );
  };

  const handleScrub = () => {
    scrubContender.mutate(undefined, {
      onSuccess: gotoScorecard,
      onError: () => toastError("Failed to remove your name."),
    });
  };

  const startTime = $derived(contest?.timeBegin ?? new Date(8640000000000000));
  const endTime = $derived(contest?.timeEnd ?? new Date(-8640000000000000));
  const gracePeriodEndTime = $derived(
    add(endTime, {
      minutes: (contest?.gracePeriod ?? 0) / (1_000_000_000 * 60),
    }),
  );
</script>

{#if !contender || !contest || !startTime || !endTime}
  <Loading />
{:else}
  <ContestStateProvider {startTime} {endTime} {gracePeriodEndTime}>
    {#snippet children({ contestState })}
      {@const disabled = contestState === "ENDED"}

      <RegistrationForm
        submit={handleSubmit}
        nameRetentionTime={contest.nameRetentionTime}
        data={{
          name: contender.name,
          compClassId: contender.compClassId,
          withdrawnFromFinals: contender.withdrawnFromFinals,
        }}
        {contestState}
      >
        <div class="controls">
          {#if !contender.scrubbedAt}
            <wa-button
              class="scrub-button"
              size="small"
              type="button"
              variant="danger"
              appearance="outlined"
              onclick={() => {
                if (scrubDialog) {
                  scrubDialog.open = true;
                }
              }}
            >
              <wa-icon slot="start" name="user-slash"></wa-icon>
              Remove my name
            </wa-button>
          {/if}
          <wa-button
            size="small"
            type="button"
            onclick={gotoScorecard}
            appearance="plain"
            >Cancel
          </wa-button>
          <wa-button
            size="small"
            type="submit"
            loading={patchContender.isPending}
            variant="neutral"
            appearance="accent"
            {disabled}
            >Save
          </wa-button>
        </div>
      </RegistrationForm>
    {/snippet}
  </ContestStateProvider>

  <wa-dialog bind:this={scrubDialog} label="Remove your name">
    Your name will be permanently removed and your results will be anonymized.
    <br /><br />
    Be aware that without a name, you will lose your chance at the finals and any
    chance of winning a raffle prize.
    <wa-button
      size="small"
      slot="footer"
      appearance="plain"
      onclick={() => {
        if (scrubDialog) {
          scrubDialog.open = false;
        }
      }}>Cancel</wa-button
    >
    <wa-button
      size="small"
      slot="footer"
      variant="danger"
      loading={scrubContender.isPending}
      onclick={handleScrub}
    >
      <wa-icon slot="start" name="user-slash"></wa-icon>
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
