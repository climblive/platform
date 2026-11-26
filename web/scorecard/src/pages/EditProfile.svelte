<script lang="ts">
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import type { ScorecardSession } from "@/types";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import type WaSwitch from "@awesome.me/webawesome/dist/components/switch/switch.js";
  import {
    CompetitionStateProvider,
    SaveIndicator,
  } from "@climblive/lib/components";
  import { checked, GenericForm, name, value } from "@climblive/lib/forms";
  import type { CompetitorPatch } from "@climblive/lib/models";
  import {
    getCompetitorQuery,
    getCompetitionQuery,
    getCategoriesQuery,
    patchCompetitorMutation,
    scrubCompetitorMutation,
  } from "@climblive/lib/queries";
  import { SyncedTime, toastUnexpectedError, z } from "@climblive/lib/utils";
  import { formatDistance, isAfter, isBefore } from "date-fns";
  import { getContext, onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Readable } from "svelte/store";
  import Loading from "./Loading.svelte";

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  const competitorQuery = $derived(getCompetitorQuery($session.competitorId));
  const categoriesQuery = $derived(getCategoriesQuery($session.competitionId));
  const patchCategory = $derived(
    patchCompetitorMutation($session.competitorId),
  );
  const patchSettings = $derived(
    patchCompetitorMutation($session.competitorId),
  );
  const profileSchema = z.object({ name: z.string().min(1) });
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

  const handleCategoryChange = (event: Event) => {
    const select = event.currentTarget as WaSelect;
    const previousValue = competitor?.categoryId?.toString() ?? "";
    const categoryId = Number(select.value);
    const category = categoriesQuery.data?.find(
      (category) => category.id === categoryId,
    );

    if (
      !category ||
      isAfter(new Date(), category.timeEnd) ||
      patchCategory.isPending
    ) {
      select.value = previousValue;
      return;
    }

    patchCategory.mutate(
      { categoryId },
      {
        onError: () => {
          select.value = previousValue;
          toastUnexpectedError("Failed to save category.");
        },
      },
    );
  };

  const handleSettingsChange = (event: Event) => {
    const toggle = event.currentTarget as WaSwitch;
    const previousValue = competitor?.withdrawnFromFinals ?? false;

    if (patchSettings.isPending) {
      toggle.checked = previousValue;
      return;
    }

    patchSettings.mutate(
      { withdrawnFromFinals: toggle.checked },
      {
        onError: () => {
          toggle.checked = previousValue;
          toastUnexpectedError("Failed to save settings.");
        },
      },
    );
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

      <div class="profile-page">
        <wa-button
          class="back-button"
          size="s"
          appearance="plain"
          onclick={gotoScorecard}
        >
          <wa-icon slot="start" name="arrow-left"></wa-icon>
          Back to scorecard
        </wa-button>

        <section aria-labelledby="profile-heading">
          <h2 id="profile-heading">Profile</h2>
          <GenericForm schema={profileSchema} submit={handleSubmit}>
            <fieldset>
              <wa-input
                size="s"
                {@attach name("name")}
                label="Name"
                type="text"
                required
                value={competitor.name}
                {disabled}
              ></wa-input>
              {@render profileCallout()}
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
            </fieldset>
          </GenericForm>
        </section>

        <section aria-labelledby="category-heading">
          <div class="section-heading">
            <h2 id="category-heading">Category</h2>
            <div class="save-status">
              {#if patchCategory.isSuccess}
                <SaveIndicator />
              {/if}
            </div>
          </div>
          <fieldset>
            <wa-select
              size="s"
              label="Category"
              required
              {@attach value(competitor.categoryId)}
              disabled={disabled ||
                patchCategory.isPending ||
                !categoriesQuery.data}
              onchange={handleCategoryChange}
            >
              {#each categoriesQuery.data ?? [] as category (category.id)}
                <wa-option
                  value={category.id}
                  disabled={isAfter(new Date(), category.timeEnd)}
                  label={category.name}
                >
                  {category.name}
                  {#if category.description}
                    <small>{category.description}</small>
                  {/if}
                </wa-option>
              {/each}
            </wa-select>
          </fieldset>
        </section>

        <section aria-labelledby="settings-heading">
          <div class="section-heading">
            <h2 id="settings-heading">Settings</h2>
            <div class="save-status">
              {#if patchSettings.isSuccess}
                <SaveIndicator />
              {/if}
            </div>
          </div>
          <fieldset>
            <wa-switch
              size="s"
              hint="If you do not wish to participate in the finals, you can give up your spot."
              disabled={disabled || patchSettings.isPending}
              {@attach checked(competitor.withdrawnFromFinals)}
              onchange={handleSettingsChange}>Opt out of finals</wa-switch
            >
          </fieldset>
        </section>
      </div>

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
  .profile-page {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-l);
    padding: var(--wa-space-m);
  }

  .back-button {
    align-self: start;

    &::part(base) {
      padding-inline: 0;
    }
  }

  section {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }

  .section-heading {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--wa-space-s);
  }

  h2 {
    margin: 0;
    font-size: var(--wa-font-size-l);
  }

  .save-status {
    display: flex;
    align-items: center;
    justify-content: end;
    min-inline-size: 5em;
  }

  fieldset {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
    margin: 0;
    padding: 0;
    min-inline-size: 0;
  }

  .controls {
    display: flex;
    flex-wrap: wrap;
    justify-content: end;
    gap: var(--wa-space-xs);

    & > .scrub-button {
      margin-right: auto;
    }
  }
</style>
