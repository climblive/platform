<script lang="ts">
  import QrCode from "@/components/QrCode.svelte";
  import TickList from "@/components/TickList.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/copy-button/copy-button.js";
  import "@awesome.me/webawesome/dist/components/divider/divider.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import WaSwitch from "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { LabeledText } from "@climblive/lib/components";
  import { checked } from "@climblive/lib/forms";
  import {
    getCompClassesQuery,
    getContenderQuery,
    patchContenderMutation,
  } from "@climblive/lib/queries";
  import { format } from "date-fns";
  import { navigate } from "svelte-routing";

  interface Props {
    contenderId: number;
  }

  let { contenderId }: Props = $props();

  const contenderQuery = $derived(getContenderQuery(contenderId));
  const patchContender = $derived(patchContenderMutation(contenderId));
  const contender = $derived($contenderQuery.data);
  const contenderName = $derived(contender?.name ?? "Unregistered");
  const contestId = $derived(contender?.contestId ?? 0);

  const compClassesQuery = $derived(
    getCompClassesQuery(contestId, { enabled: !!contestId }),
  );

  const compClasses = $derived($compClassesQuery.data);

  let withdrawFromFinalsToggle: WaSwitch | undefined = $state();

  const handleToggleWithdrawFromFinals = () => {
    if (!withdrawFromFinalsToggle) {
      return;
    }

    $patchContender.mutate({
      withdrawnFromFinals: withdrawFromFinalsToggle.checked,
    });
  };

  const handleDisqualify = () => {
    $patchContender.mutate({
      disqualified: true,
    });
  };

  const handleRequalify = () => {
    $patchContender.mutate({
      disqualified: false,
    });
  };
</script>

{#if contender && compClasses}
  <wa-button
    appearance="plain"
    onclick={() => navigate(`/admin/contests/${contestId}#results`)}
    >Back to results<wa-icon name="arrow-left" slot="start"
    ></wa-icon></wa-button
  >

  <h1>
    {#if contender.disqualified}
      <strike>
        {contenderName}
      </strike>
    {:else}
      {contenderName}
    {/if}
  </h1>
  <section>
    <article>
      <LabeledText label="Class"
        >{compClasses.find(({ id }) => id === contender.compClassId)
          ?.name}</LabeledText
      >
      <LabeledText label="Club">{contender.clubName}</LabeledText>
      {#if contender.entered}
        <LabeledText label="Entered"
          >{format(contender.entered, "yyyy-MM-dd HH:mm")}</LabeledText
        >
      {/if}
      {#if contender.disqualified}
        <LabeledText label="Disqualified">Yes</LabeledText>
      {/if}
      {#if !contender.disqualified}
        <LabeledText label="Placement"
          >{contender.score?.placement ?? "-"}</LabeledText
        >
        <LabeledText label="Score">{contender.score?.score ?? "-"}</LabeledText>
        <LabeledText label="Finalist">
          <wa-icon name={contender.score?.finalist ? "medal" : "minus"}
          ></wa-icon>
        </LabeledText>
      {/if}

      <LabeledText label="Registration code">
        {contender.registrationCode}
        <wa-copy-button
          value={`${location.protocol}//${location.host}/${contender.registrationCode}`}
        ></wa-copy-button>
      </LabeledText>
    </article>
    <div class="registration">
      <QrCode
        registrationCode={contender.registrationCode}
        width={200}
        fill={contender.disqualified
          ? "var(--wa-color-text-quiet)"
          : "var(--wa-color-text-normal)"}
      ></QrCode>

      <wa-button
        href={`/${contender.registrationCode}`}
        target="_blank"
        appearance="plain"
        variant="brand"
        size="large"
      >
        <wa-icon slot="start" name="arrow-up-right-from-square"></wa-icon>
        Open scorecard
      </wa-button>
    </div>
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
  <wa-divider style="--color: var(--wa-color-brand-fill-normal);"></wa-divider>
  <div class="controls">
    <wa-switch
      bind:this={withdrawFromFinalsToggle}
      hint="The contender does not wish to take part in the finals"
      {@attach checked(contender.withdrawnFromFinals)}
      onchange={handleToggleWithdrawFromFinals}
      disabled={contender.disqualified || $patchContender.isPending}
      >Withdraw from finals</wa-switch
    >

    {#if contender.disqualified}
      <wa-button
        onclick={handleRequalify}
        loading={$patchContender.isPending}
        appearance="outlined"
        >Requalify contender<wa-icon slot="start" name="right-to-bracket"
        ></wa-icon></wa-button
      >
    {:else}
      <wa-button
        variant="danger"
        onclick={handleDisqualify}
        loading={$patchContender.isPending}
        >Disqualify contender<wa-icon slot="start" name="skull-crossbones"
        ></wa-icon></wa-button
      >
    {/if}
  </div>
{/if}

<style>
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
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }

  .registration {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
    align-items: center;

    & wa-button::part(base) {
      padding-inline: 0;
    }
  }

  .copy {
    color: var(--wa-color-text-quiet);
  }

  .controls {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-l);
  }
</style>
