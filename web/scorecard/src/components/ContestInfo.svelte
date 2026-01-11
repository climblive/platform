<script lang="ts">
  import "@awesome.me/webawesome/dist/components/details/details.js";
  import type WaDetails from "@awesome.me/webawesome/dist/components/details/details.js";
  import { LabeledText } from "@climblive/lib/components";
  import type { CompClass, Contest, Problem } from "@climblive/lib/models";
  import { format } from "date-fns";
  import { sv } from "date-fns/locale";

  interface Props {
    contest: Contest;
    compClasses: CompClass[];
    problems: Problem[];
  }

  let { contest, compClasses, problems }: Props = $props();

  let details: WaDetails | undefined = $state();

  const scoreboardUrl = `${location.protocol}//${location.host}/scoreboard/${contest.id}`;
</script>

<section>
  <LabeledText label="Name">{contest.name}</LabeledText>
  {#if contest.description}
    <LabeledText label="Description">{contest.description}</LabeledText>
  {/if}
  {#if contest.location}
    <LabeledText label="Location">{contest.location}</LabeledText>
  {/if}
  {#if contest.timeBegin}
    <LabeledText label="Start time">
      {format(contest.timeBegin, "PPpp", { locale: sv })}
    </LabeledText>
  {/if}
  {#if contest.timeEnd}
    <LabeledText label="End time">
      {format(contest.timeEnd, "PPpp", { locale: sv })}
    </LabeledText>
  {/if}
  <LabeledText label="Competition classes">
    {compClasses.map((cc) => cc.name).join(", ")}
  </LabeledText>
  <LabeledText label="Number of problems">
    {problems.length.toString()}</LabeledText
  >
  <LabeledText label="Qualifying problems">
    {`${contest.qualifyingProblems.toString()} hardest`}
  </LabeledText>
  <LabeledText label="Number of finalists">
    {contest.finalists.toString()}
  </LabeledText>
  <LabeledText label="Scoreboard">
    <a href={scoreboardUrl} target="_blank">{scoreboardUrl}</a>
  </LabeledText>
</section>
{#if contest.info}
  <wa-details
    onwa-after-show={() =>
      details?.scrollIntoView({
        behavior: "smooth",
        block: "start",
        inline: "nearest",
      })}
    bind:this={details}
    summary="General info"
  >
    {@html contest.info}
  </wa-details>
{/if}

<style>
  section {
    padding: var(--wa-space-m);
    background-color: var(--wa-color-surface-default);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-surface-border);
    border-radius: var(--wa-border-radius-m);
    font-size: var(--wa-font-size-s);

    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }

  wa-details::part(base) {
    margin-top: var(--wa-space-s);
    font-size: var(--wa-font-size-s);
    border-radius: var(--wa-border-radius-m);
  }

  wa-details::part(content) {
    padding-top: 0;
  }
</style>
