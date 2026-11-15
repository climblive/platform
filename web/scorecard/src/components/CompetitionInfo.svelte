<script lang="ts">
  import "@awesome.me/webawesome/dist/components/details/details.js";
  import type WaDetails from "@awesome.me/webawesome/dist/components/details/details.js";
  import { LabeledText } from "@climblive/lib/components";
  import type { CompClass, Competition, Problem } from "@climblive/lib/models";
  import { format } from "date-fns";
  import { sv } from "date-fns/locale";

  interface Props {
    competition: Competition;
    compClasses: CompClass[];
    problems: Problem[];
  }

  let { competition, compClasses, problems }: Props = $props();

  let details: WaDetails | undefined = $state();

  const scoreboardUrl = `${location.protocol}//${location.host}/scoreboard/${competition.id}`;
</script>

<section>
  <LabeledText label="Competition">{competition.name}</LabeledText>
  {#if competition.description}
    <LabeledText label="Description">{competition.description}</LabeledText>
  {/if}
  {#if competition.location}
    <LabeledText label="Location">{competition.location}</LabeledText>
  {/if}
  {#if competition.timeBegin}
    <LabeledText label="Start time">
      {format(competition.timeBegin, "PPpp", { locale: sv })}
    </LabeledText>
  {/if}
  {#if competition.timeEnd}
    <LabeledText label="End time">
      {format(competition.timeEnd, "PPpp", { locale: sv })}
    </LabeledText>
  {/if}
  <LabeledText label="Categories">
    {compClasses.map((cc) => cc.name).join(", ")}
  </LabeledText>
  <LabeledText label="Number of problems">
    {problems.length.toString()}</LabeledText
  >
  {#if competition.qualifyingProblems > 0}
    <LabeledText label="Qualifying problems">
      {`${competition.qualifyingProblems.toString()} hardest`}
    </LabeledText>
  {/if}
  {#if competition.finalists > 0}
    <LabeledText label="Number of finalists">
      {competition.finalists.toString()}
    </LabeledText>
  {/if}
  <LabeledText label="Scoreboard">
    <a href={scoreboardUrl} target="_blank">{scoreboardUrl}</a>
  </LabeledText>
</section>
{#if competition.info}
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
    {@html competition.info}
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
