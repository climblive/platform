<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/details/details.js";
  import type WaDetails from "@awesome.me/webawesome/dist/components/details/details.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { LabeledText } from "@climblive/lib/components";
  import type { Category, Competition, Problem } from "@climblive/lib/models";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import { format } from "date-fns";
  import { sv } from "date-fns/locale";

  interface Props {
    competition: Competition;
    categories: Category[];
    problems: Problem[];
  }

  let { competition, categories, problems }: Props = $props();

  let details: WaDetails | undefined = $state();
  let sharing = $state(false);

  const scoreboardUrl = $derived(
    `${location.origin}/scoreboard/${competition.id}`,
  );

  async function shareScoreboard() {
    sharing = true;

    try {
      await navigator.share({
        title: competition.name,
        url: scoreboardUrl,
      });
    } catch (error) {
      if (!(error instanceof DOMException && error.name === "AbortError")) {
        toastUnexpectedError("Failed to share scoreboard.");
      }
    } finally {
      sharing = false;
    }
  }
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
    {categories.map((cc) => cc.name).join(", ")}
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
    <span class="scoreboard-link">
      <a href={scoreboardUrl} target="_blank">{scoreboardUrl}</a>
      {#if typeof navigator.share === "function"}
        <wa-button
          size="small"
          appearance="plain"
          disabled={sharing}
          onclick={shareScoreboard}
        >
          <wa-icon name="share-nodes" label="Share scoreboard"></wa-icon>
        </wa-button>
      {/if}
    </span>
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

  .scoreboard-link {
    display: flex;
    align-items: center;
    gap: var(--wa-space-xs);

    a {
      min-width: 0;
      overflow-wrap: anywhere;
    }

    wa-button {
      flex-shrink: 0;
    }
  }
</style>
