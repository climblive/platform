<script lang="ts">
  import type { CompClass, Contest, Problem } from "@climblive/lib/models";
  import type { SlDetails } from "@shoelace-style/shoelace";
  import { format } from "date-fns";
  import { sv } from "date-fns/locale";
  import LabeledText from "./LabeledText.svelte";

  export let contest: Contest;
  export let compClasses: CompClass[];
  export let problems: Problem[];

  let details: SlDetails;

  $: {
    if (details && contest.rules) {
      details.innerHTML = contest.rules;
    }
  }

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
      {format(contest.timeEnd, "PPpp", { locale: sv })}}
    </LabeledText>
  {/if}
  <LabeledText label="Classes">
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
{#if contest.rules}
  <sl-details bind:this={details} summary="Rules"> </sl-details>
{/if}

<style>
  section {
    padding: var(--sl-spacing-medium);
    color: var(--sl-color-neutral-700);
    background-color: var(--sl-color-primary-100);
    border: solid 1px
      color-mix(in srgb, var(--sl-color-primary-300), transparent 50%);
    border-radius: var(--sl-border-radius-small);
    font-size: var(--sl-font-size-small);

    display: flex;
    flex-direction: column;
    gap: var(--sl-spacing-medium);
  }

  sl-details::part(base) {
    margin-top: var(--sl-spacing-small);
    background-color: var(--sl-color-primary-100);
    font-size: var(--sl-font-size-small);
  }

  sl-details::part(content) {
    padding-top: 0;
  }

  a {
    color: var(--sl-color-primary-700);
  }
</style>
