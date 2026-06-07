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
    ContenderName,
    EmptyState,
    LabeledText,
  } from "@climblive/lib/components";
  import { checked, value } from "@climblive/lib/forms";
  import {
    getCompClassesQuery,
    getContenderQuery,
    getContestQuery,
    patchContenderMutation,
  } from "@climblive/lib/queries";
  import { ordinalSuperscript } from "@climblive/lib/utils";
  import { format } from "date-fns";
  import { navigate } from "svelte-routing";

  interface Props {
    contenderId: number;
  }

  let { contenderId }: Props = $props();

  const contenderQuery = $derived(getContenderQuery(contenderId));
  const patchContender = $derived(patchContenderMutation(contenderId));
  const contender = $derived(contenderQuery.data);
  const contestId = $derived(contender?.contestId ?? 0);

  const compClassesQuery = $derived(
    getCompClassesQuery(contestId, { enabled: !!contestId }),
  );

  const compClasses = $derived(compClassesQuery.data);

  const contestQuery = $derived(
    contestId ? getContestQuery(contestId) : undefined,
  );
  const contest = $derived(contestQuery?.data);

  let withdrawFromFinalsToggle: WaSwitch | undefined = $state();
  let compClassSelect: WaSelect | undefined = $state();

  const handleToggleWithdrawFromFinals = () => {
    if (!withdrawFromFinalsToggle) {
      return;
    }

    patchContender.mutate({
      withdrawnFromFinals: withdrawFromFinalsToggle.checked,
    });
  };

  const handleDisqualify = () => {
    patchContender.mutate({
      disqualified: true,
    });
  };

  const handleRequalify = () => {
    patchContender.mutate({
      disqualified: false,
    });
  };

  const handleCompClassChange = () => {
    if (!compClassSelect || !compClassSelect.value) {
      return;
    }

    if (typeof compClassSelect.value !== "string") {
      return;
    }

    const newCompClassId = parseInt(compClassSelect.value, 10);

    patchContender.mutate({
      compClassId: newCompClassId,
    });
  };
</script>

{#if contender && compClasses && contest}
  <wa-breadcrumb>
    <wa-breadcrumb-item
      onclick={() =>
        navigate(`/admin/organizers/${contest.ownership.organizerId}/contests`)}
      ><wa-icon name="home"></wa-icon></wa-breadcrumb-item
    >
    <wa-breadcrumb-item onclick={() => navigate(`/admin/contests/${contestId}`)}
      >{contest.name}</wa-breadcrumb-item
    >
    <wa-breadcrumb-item
      onclick={() => navigate(`/admin/contests/${contestId}/results`)}
      >Results</wa-breadcrumb-item
    >
  </wa-breadcrumb>

  <h1>
    {#if contender.entered}
      {#if contender.disqualified}
        <del>
          <ContenderName
            id={contender.id}
            name={contender.name}
            scrubbedAt={contender.scrubbedAt}
            withTooltip
          />
        </del>
      {:else}
        <ContenderName
          id={contender.id}
          name={contender.name}
          scrubbedAt={contender.scrubbedAt}
          withTooltip
        />
      {/if}
    {/if}
  </h1>

  {#if !contender.entered}
    <EmptyState
      title="Not registered"
      description="Have the contender scan the QR code below to enter the contest and start climbing."
    >
      {#snippet actions()}
        <QrCode
          registrationCode={contender.registrationCode}
          width={200}
          fill="var(--wa-color-text-normal)"
        ></QrCode>

        <wa-button
          href={`/${contender.registrationCode}`}
          target="_blank"
          appearance="plain"
          variant="neutral"
          size="l"
        >
          <wa-icon slot="start" name="arrow-up-right-from-square"></wa-icon>
          Open scorecard
        </wa-button>

        <span class="registration-code">
          {contender.registrationCode}
          <wa-copy-button value={contender.registrationCode}></wa-copy-button>
        </span>
      {/snippet}
    </EmptyState>
  {:else}
    <div class="codes">
      <a
        href={`/${contender.registrationCode}`}
        target="_blank"
        class="open-scorecard-link"
      >
        <wa-icon name="arrow-up-right-from-square"></wa-icon>
        Open scorecard
      </a>

      <span class="registration-code">
        {contender.registrationCode}
        <wa-copy-button value={contender.registrationCode} size="l"
        ></wa-copy-button>
      </span>
    </div>

    <section>
      <article>
        <LabeledText label="Class"
          >{compClasses.find(({ id }) => id === contender.compClassId)?.name ??
            "-"}</LabeledText
        >
        <span class="entered">
          <LabeledText label="Entered"
            >{format(contender.entered, "yyyy-MM-dd HH:mm")}</LabeledText
          ></span
        >
        <LabeledText label="Placement">
          {#if contender.disqualified}
            Disqualified
          {:else if contender.score?.placement}
            {contender.score?.placement}<sup
              >{ordinalSuperscript(contender.score.placement)}</sup
            >
          {:else}
            -
          {/if}
        </LabeledText>
        <LabeledText label="Score">{contender.score?.score ?? "-"}</LabeledText>
        <LabeledText label="Finalist">
          <wa-icon name={contender.score?.finalist ? "medal" : "minus"}
          ></wa-icon>
        </LabeledText>
      </article>
    </section>

    {#if !contender.disqualified}
      <h2>Ticks</h2>
      <wa-divider style="--color: var(--wa-color-brand-fill-normal);"
      ></wa-divider>
      <p class="copy">
        All ticks registered by the contender during the contest.
      </p>
      <TickList contenderId={contender.id} contestId={contender.contestId}
      ></TickList>
    {/if}

    <h2>Settings</h2>
    <wa-divider style="--color: var(--wa-color-brand-fill-normal);"
    ></wa-divider>
    <div class="controls">
      <wa-select
        size="s"
        bind:this={compClassSelect}
        label="Competition class"
        hint="Change the class for this contender."
        {@attach value(contender.compClassId)}
        onchange={handleCompClassChange}
        disabled={contender.disqualified || patchContender.isPending}
      >
        {#each compClasses as compClass (compClass.id)}
          <wa-option value={compClass.id} label={compClass.name}>
            {compClass.name}
            {#if compClass.description}
              <small>{compClass.description}</small>
            {/if}
          </wa-option>
        {/each}
      </wa-select>

      <wa-switch
        size="s"
        bind:this={withdrawFromFinalsToggle}
        hint="In case the contender does not wish to take part in the finals."
        {@attach checked(contender.withdrawnFromFinals)}
        onchange={handleToggleWithdrawFromFinals}
        disabled={contender.disqualified || patchContender.isPending}
        >Withdraw from finals</wa-switch
      >

      {#if contender.disqualified}
        <wa-button
          onclick={handleRequalify}
          loading={patchContender.isPending}
          appearance="outlined"
          >Requalify contender<wa-icon slot="start" name="right-to-bracket"
          ></wa-icon></wa-button
        >
      {:else}
        <wa-button
          size="s"
          variant="danger"
          onclick={handleDisqualify}
          loading={patchContender.isPending}
          >Disqualify contender<wa-icon slot="start" name="user-slash"
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
    font-size: var(--wa-font-size-s);
  }

  .codes {
    margin-block-end: var(--wa-space-m);
    display: flex;
    gap: var(--wa-space-m);
    justify-content: space-between;
    align-items: center;
  }
</style>
