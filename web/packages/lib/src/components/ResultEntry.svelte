<script lang="ts">
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import type { ScoreboardEntry } from "../models";
  import { ordinalSuperscript } from "../utils";
  import Score from "./Score.svelte";

  interface Props {
    scoreboardEntry: ScoreboardEntry;
    highlighted?: boolean;
  }

  let { scoreboardEntry, highlighted = false }: Props = $props();
  let score = $derived(scoreboardEntry.score);
</script>

<section
  data-finalist={score?.finalist ?? false}
  data-highlighted={highlighted ? "true" : "false"}
>
  <div class="number">
    {#if score?.placement}
      {score.placement}<sup>{ordinalSuperscript(score.placement)}</sup>
    {:else}
      -
    {/if}
  </div>
  <div class="name">{scoreboardEntry.name}</div>
  <div class="score">
    {#if score === undefined || score.score === 0}
      -
    {:else}
      <Score value={score.score} />
      <wa-icon name={score?.finalist ? "medal" : "minus"}></wa-icon>
    {/if}
  </div>
</section>

<style>
  section {
    height: 2.25rem;
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
    background-color: var(--wa-color-primary-fill-quiet);
    animation: highlight-fade 2s ease-in-out forwards;
  }

  @keyframes highlight-fade {
    0% {
      background-color: var(--wa-color-primary-fill-normal);
    }
    100% {
      background-color: var(--wa-color-primary-fill-quiet);
    }
  }

  section > div {
    width: 100%;
    white-space: nowrap;
    overflow: hidden;
    line-height: 2.25rem;
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
