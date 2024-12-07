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
</script>

<section>
  <LabeledText label="Name" text={contest.name} />
  {#if contest.description}
    <LabeledText label="Description" text={contest.description} />
  {/if}
  {#if contest.location}
    <LabeledText label="Location" text={contest.location} />
  {/if}
  {#if contest.timeBegin}
    <LabeledText
      label="Start time"
      text={format(contest.timeBegin, "PPpp", { locale: sv })}
    />
  {/if}
  {#if contest.timeEnd}
    <LabeledText
      label="End time"
      text={format(contest.timeEnd, "PPpp", { locale: sv })}
    />
  {/if}
  <LabeledText
    label="Classes"
    text={compClasses.map((cc) => cc.name).join(", ")}
  />
  <LabeledText label="Number of problems" text={problems.length.toString()} />
  <LabeledText
    label="Qualifying problems"
    text={`${contest.qualifyingProblems.toString()} hardest`}
  />
  <LabeledText
    label="Number of finalists"
    text={contest.finalists.toString()}
  />
  {#if contest.rules}
    <sl-details bind:this={details} summary="Rules"> </sl-details>
  {/if}
</section>

<style>
  section {
    padding-inline: var(--sl-spacing-small);
    color: var(--sl-color-neutral-700);
  }

  sl-details::part(base) {
    background-color: transparent;
    border-color: var(--sl-color-primary-600);
  }
</style>
