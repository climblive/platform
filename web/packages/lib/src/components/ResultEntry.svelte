<script lang="ts">
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import type { ScoreboardEntry } from "../models";
  import { NO_SCORE } from "../models/score";
  import { ordinalSuperscript } from "../utils";
  import CompetitorName from "./CompetitorName.svelte";
  import Score from "./Score.svelte";

  interface Props {
    scoreboardEntry: ScoreboardEntry;
    highlighted?: boolean;
  }

  const { scoreboardEntry, highlighted = false }: Props = $props();
  let result = $derived(scoreboardEntry.result);
</script>

<section
  data-finalist={result?.finalist ?? false}
  data-highlighted={highlighted ? "true" : "false"}
>
  <div class="number">
    {#if result?.placement}
      {result.placement}<sup>{ordinalSuperscript(result.placement)}</sup>
    {:else}
      -
    {/if}
  </div>
  <div class="name">
    <CompetitorName
      id={scoreboardEntry.competitorId}
      name={scoreboardEntry.name}
      scrubbedAt={scoreboardEntry.scrubbedAt}
    />
  </div>
  <div class="score">
    {#if result === undefined || result.score === NO_SCORE}
      -
    {:else}
      <Score value={result.score} animate />
      <wa-icon name={result?.finalist ? "medal" : "minus"}></wa-icon>
    {/if}
  </div>
</section>

<style>
  section {
    height: 2.75rem;
    background-color: var(--wa-color-surface-raised);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-neutral-border-quiet);
    display: grid;
    padding-inline: var(--wa-space-s);
    border-radius: var(--wa-border-radius-m);
    grid-template-columns: 2rem 1fr max-content;
    grid-template-rows: 1fr;
    gap: var(--wa-space-xs);
    align-items: center;
    justify-items: start;
    font-size: var(--wa-font-size-m);
    font-weight: var(--wa-font-weight-semibold);
    user-select: none;
  }

  section[data-highlighted="true"] {
    border-color: var(--wa-color-brand-border-loud);
  }

  section > div {
    width: 100%;
    white-space: nowrap;
    overflow: hidden;
    line-height: 2.75rem;
  }

  .number {
    font-size: var(--wa-font-size-xs);
  }

  .name {
    text-overflow: ellipsis;
  }

  .score {
    justify-self: end;

    display: flex;
    align-items: center;
    gap: var(--wa-space-xs);

    text-align: right;
    font-weight: var(--wa-font-weight-bold);
    font-size: var(--wa-font-size-m);
  }
</style>
