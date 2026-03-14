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
  import { SyncedTime, toastError } from "@climblive/lib/utils";
  import { formatDistance } from "date-fns";
  import { getContext, onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Readable } from "svelte/store";
  import Loading from "./Loading.svelte";

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  const contenderQuery = $derived(getContenderQuery($session.contenderId));
  const contestQuery = $derived(getContestQuery($session.contestId));
  const patchContender = $derived(patchContenderMutation($session.contenderId));
  const scrubContender = $derived(scrubContenderMutation($session.contenderId));
  const time = new SyncedTime(60_000);

  onMount(() => {
    time.start();

    return () => time.stop();
  });

  let contender = $derived(contenderQuery.data);
  let contest = $derived(contestQuery.data);

  const retentionDuration = $derived.by(() => {
    if (!contender?.scrubBefore) {
      return undefined;
    }

    const now = time.current;

    if (contender.scrubBefore <= now) {
      return undefined;
    }

    return formatDistance(contender.scrubBefore, now);
  });

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
</script>

{#if !contender || !contest}
  <Loading />
{:else}
  <ContestStateProvider
    contestId={contest.id}
    compClassId={contender.compClassId}
  >
    {#snippet children({ contestState })}
      {@const disabled = contestState === "ENDED"}

      <RegistrationForm
        submit={handleSubmit}
        data={{
          name: contender.name,
          compClassId: contender.compClassId,
          withdrawnFromFinals: contender.withdrawnFromFinals,
        }}
        callout={profileCallout}
        {contestState}
      >
        <div class="controls">
          {#if !!contender.name}
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

      {#snippet profileCallout()}
        <wa-callout variant="neutral" size="small">
          <wa-icon slot="icon" name="circle-info"></wa-icon>
          {#if retentionDuration}
            Your name will be kept stored for {retentionDuration} from now, after
            which it will be removed and your results anonymized.
          {:else}
            Your name will be removed and your results anonymized shortly.
          {/if}
        </wa-callout>
      {/snippet}
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
