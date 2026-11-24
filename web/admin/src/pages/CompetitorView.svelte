<script lang="ts">
  import QrCode from "@/components/QrCode.svelte";
  import TickList from "@/components/TickList.svelte";
  import "@awesome.me/webawesome/dist/components/breadcrumb-item/breadcrumb-item.js";
  import "@awesome.me/webawesome/dist/components/breadcrumb/breadcrumb.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/copy-button/copy-button.js";
  import "@awesome.me/webawesome/dist/components/divider/divider.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import WaSwitch from "@awesome.me/webawesome/dist/components/switch/switch.js";
  import {
    CompetitorName,
    EmptyState,
    LabeledText,
    RegistrationCode,
    Score,
  } from "@climblive/lib/components";
  import { checked, value } from "@climblive/lib/forms";
  import {
    getCategoriesQuery,
    getCompetitorQuery,
    getCompetitionQuery,
    patchCompetitorMutation,
  } from "@climblive/lib/queries";
  import {
    formatCategoryEligibility,
    ordinalSuperscript,
  } from "@climblive/lib/utils";
  import { format } from "date-fns";
  import { navigate } from "svelte-routing";

  interface Props {
    competitorId: number;
  }

  let { competitorId }: Props = $props();

  const competitorQuery = $derived(getCompetitorQuery(competitorId));
  const patchCompetitor = $derived(patchCompetitorMutation(competitorId));
  const competitor = $derived(competitorQuery.data);
  const competitionId = $derived(competitor?.competitionId ?? 0);

  const categoriesQuery = $derived(
    getCategoriesQuery(competitionId, { enabled: !!competitionId }),
  );

  const categories = $derived(categoriesQuery.data);

  const competitionQuery = $derived(
    competitionId ? getCompetitionQuery(competitionId) : undefined,
  );
  const competition = $derived(competitionQuery?.data);

  let withdrawFromFinalsToggle: WaSwitch | undefined = $state();
  let categorySelect: WaSelect | undefined = $state();

  const handleToggleWithdrawFromFinals = () => {
    if (!withdrawFromFinalsToggle) {
      return;
    }

    patchCompetitor.mutate({
      withdrawnFromFinals: withdrawFromFinalsToggle.checked,
    });
  };

  const handleDisqualify = () => {
    patchCompetitor.mutate({
      disqualified: true,
    });
  };

  const handleRequalify = () => {
    patchCompetitor.mutate({
      disqualified: false,
    });
  };

  const handleCategoryChange = () => {
    if (!categorySelect || !categorySelect.value) {
      return;
    }

    if (typeof categorySelect.value !== "string") {
      return;
    }

    const newCategoryId = parseInt(categorySelect.value, 10);

    patchCompetitor.mutate({
      categoryId: newCategoryId,
    });
  };
</script>

{#if competitor && categories && competition}
  <wa-breadcrumb>
    <wa-breadcrumb-item
      onclick={() =>
        navigate(
          `/admin/organizers/${competition.ownership.organizerId}/competitions`,
        )}><wa-icon name="home"></wa-icon></wa-breadcrumb-item
    >
    <wa-breadcrumb-item
      onclick={() => navigate(`/admin/competitions/${competitionId}`)}
      >{competition.name}</wa-breadcrumb-item
    >
    <wa-breadcrumb-item
      onclick={() => navigate(`/admin/competitions/${competitionId}/results`)}
      >Results</wa-breadcrumb-item
    >
  </wa-breadcrumb>

  <h1>
    {#if competitor.entered}
      {#if competitor.disqualified}
        <del>
          <CompetitorName
            id={competitor.id}
            name={competitor.name}
            scrubbedAt={competitor.scrubbedAt}
            withTooltip
          />
        </del>
      {:else}
        <CompetitorName
          id={competitor.id}
          name={competitor.name}
          scrubbedAt={competitor.scrubbedAt}
          withTooltip
        />
      {/if}
    {/if}
  </h1>

  {#if !competitor.entered}
    <EmptyState
      title="Not registered"
      description="Have a competitor scan the QR code below to register for the competition."
    >
      {#snippet actions()}
        <QrCode
          registrationCode={competitor.registrationCode}
          width={200}
          fill="var(--wa-color-text-normal)"
        ></QrCode>

        <wa-button
          href={`/${competitor.registrationCode}`}
          target="_blank"
          appearance="plain"
          variant="neutral"
          size="l"
        >
          <wa-icon slot="start" name="arrow-up-right-from-square"></wa-icon>
          Open scorecard
        </wa-button>
      {/snippet}
    </EmptyState>
  {:else}
    <div class="codes">
      <a
        href={`/${competitor.registrationCode}`}
        target="_blank"
        class="open-scorecard-link"
      >
        <wa-icon name="arrow-up-right-from-square"></wa-icon>
        Open scorecard
      </a>

      <span class="registration-code">
        <RegistrationCode code={competitor.registrationCode} />
        <wa-copy-button value={competitor.registrationCode}></wa-copy-button>
      </span>
    </div>

    <section>
      <article>
        <LabeledText label="Category"
          >{categories.find(({ id }) => id === competitor.categoryId)?.name ??
            "-"}</LabeledText
        >
        <span class="entered">
          <LabeledText label="Entered"
            >{format(competitor.entered, "yyyy-MM-dd HH:mm")}</LabeledText
          ></span
        >
        <LabeledText label="Placement">
          {#if competitor.disqualified}
            Disqualified
          {:else if competitor.result?.placement}
            {competitor.result.placement}<sup
              >{ordinalSuperscript(competitor.result.placement)}</sup
            >
          {:else}
            -
          {/if}
        </LabeledText>
        <LabeledText label="Score">
          {#if competitor.result?.score}
            <Score value={competitor.result.score}></Score>
          {:else}
            -
          {/if}
        </LabeledText>
        <LabeledText label="Finalist">
          <wa-icon name={competitor.result?.finalist ? "medal" : "minus"}
          ></wa-icon>
        </LabeledText>
      </article>
    </section>

    {#if !competitor.disqualified}
      <h2>Results</h2>
      <wa-divider style="--color: var(--wa-color-brand-fill-normal);"
      ></wa-divider>
      <p class="copy">
        Problems attempted by the competitor during the competition.
      </p>
      <TickList
        competitorId={competitor.id}
        competitionId={competitor.competitionId}
      ></TickList>
    {/if}

    <h2>Settings</h2>
    <wa-divider style="--color: var(--wa-color-brand-fill-normal);"
    ></wa-divider>
    <div class="controls">
      <wa-select
        size="s"
        bind:this={categorySelect}
        label="Category"
        hint="Change the category for this competitor."
        {@attach value(competitor.categoryId)}
        onchange={handleCategoryChange}
        disabled={competitor.disqualified || patchCompetitor.isPending}
      >
        {#each categories as category (category.id)}
          <wa-option value={category.id} label={category.name}>
            {category.name}
            <small>{formatCategoryEligibility(category)}</small>
          </wa-option>
        {/each}
      </wa-select>

      <wa-switch
        size="s"
        bind:this={withdrawFromFinalsToggle}
        hint="In case the competitor does not wish to take part in the finals."
        {@attach checked(competitor.withdrawnFromFinals)}
        onchange={handleToggleWithdrawFromFinals}
        disabled={competitor.disqualified || patchCompetitor.isPending}
        >Withdraw from finals</wa-switch
      >

      {#if competitor.disqualified}
        <wa-button
          onclick={handleRequalify}
          loading={patchCompetitor.isPending}
          appearance="outlined"
          >Requalify competitor<wa-icon slot="start" name="right-to-bracket"
          ></wa-icon></wa-button
        >
      {:else}
        <wa-button
          size="s"
          variant="danger"
          onclick={handleDisqualify}
          loading={patchCompetitor.isPending}
          >Disqualify competitor<wa-icon slot="start" name="user-slash"
          ></wa-icon></wa-button
        >
      {/if}
    </div>
  {/if}
{/if}

<style>
  wa-breadcrumb {
    margin-block-end: var(--wa-space-m);
    display: block;
  }

  h2 {
    margin-top: var(--wa-space-2xl);
  }

  section {
    display: flex;
    gap: var(--wa-space-m);
    width: 100%;
    justify-content: space-between;
    flex-wrap: wrap;
  }

  article {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    row-gap: var(--wa-space-s);
    column-gap: var(--wa-space-2xl);
    align-items: start;
    width: 100%;
  }

  .copy {
    color: var(--wa-color-text-quiet);
  }

  .controls {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-l);
  }

  .registration-code {
    display: grid;
    grid-template-columns: 1fr max-content;
    align-items: center;
    font-family: monospace;
    font-size: var(--wa-font-size-m);
  }

  .entered {
    grid-column: 2 / span 2;
  }

  .open-scorecard-link {
    display: block;
    font-size: var(--wa-font-size-m);
  }

  .codes {
    margin-block-end: var(--wa-space-m);
    display: flex;
    gap: var(--wa-space-s);
    justify-content: space-between;
    align-items: center;
  }
</style>
