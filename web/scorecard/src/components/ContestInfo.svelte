<script lang="ts">
  import type { CompClass, Contest, Problem } from "@climblive/lib/models";
  import type { SlDetails } from "@shoelace-style/shoelace";
  import "@shoelace-style/shoelace/dist/components/details/details.js";
  import { format } from "date-fns";
  import { sv } from "date-fns/locale";
  import LabeledText from "./LabeledText.svelte";

  interface Props {
    contest: Contest;
    compClasses: CompClass[];
    problems: Problem[];
  }

  let { contest, compClasses, problems }: Props = $props();

  let details: SlDetails | undefined = $state();

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
{#if contest.rules}
  <sl-details
    onsl-after-show={() =>
      details?.scrollIntoView({
        behavior: "smooth",
        block: "start",
        inline: "nearest",
      })}
    bind:this={details}
    summary="Rules"
  >
    {@html contest.rules}
  </sl-details>
{/if}

<style>
  section {
    padding: var(--sl-spacing-medium);
    background-color: var(--sl-color-neutral-50);
    border: solid 1px var(--sl-color-neutral-300);
    border-radius: var(--sl-border-radius-small);
    font-size: var(--sl-font-size-small);

    display: flex;
    flex-direction: column;
    gap: var(--sl-spacing-medium);
  }

  sl-details::part(base) {
    margin-top: var(--sl-spacing-small);
    background-color: var(--sl-color-neutral-50);
    border-color: var(--sl-color-neutral-300);
    font-size: var(--sl-font-size-small);
  }

  sl-details::part(content) {
    padding-top: 0;
  }

  a {
    color: var(--sl-color-primary-700);
  }
</style>
