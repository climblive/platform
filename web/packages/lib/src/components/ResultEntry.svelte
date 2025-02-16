<script lang="ts">
  import "@shoelace-style/shoelace/dist/components/icon/icon.js";
  import type { ScoreboardEntry } from "../models";
  import { ordinalSuperscript } from "../utils";
  import Score from "./Score.svelte";

  interface Props {
    scoreboardEntry: ScoreboardEntry;
  }

  let { scoreboardEntry }: Props = $props();
  let score = $derived(scoreboardEntry.score);
</script>

<section data-finalist={score?.finalist ?? false}>
  <div class="number">
    {#if score?.placement}
      {score.placement}<sup>{ordinalSuperscript(score.placement)}</sup>
    {:else}
      -
    {/if}
  </div>
  <div class="name">{scoreboardEntry.publicName}</div>
  <div class="score">
    {#if score === undefined || score.score === 0}
      -
    {:else}
      <Score value={score.score} />
      <sl-icon name={score?.finalist ? "trophy" : "dash"}></sl-icon>
    {/if}
  </div>
</section>

<style>
  section {
    height: 2.25rem;
    background-color: var(--sl-color-primary-500);
    color: white;
    display: grid;
    padding-inline: var(--sl-spacing-small);
    border-radius: var(--sl-border-radius-medium);
    grid-template-columns: 2rem 1fr max-content;
    grid-template-rows: 1fr;
    gap: var(--sl-spacing-x-small);
    align-items: center;
    justify-items: start;
    font-size: var(--sl-font-size-medium);
    font-weight: var(--sl-font-weight-semibold);
    user-select: none;
  }

  section[data-finalist="true"] {
    background: linear-gradient(
      45deg,
      var(--sl-color-yellow-800),
      var(--sl-color-yellow-500)
    );
    background-size: 150%;
    background-position: left;
    animation: shine 5s infinite;
    color: white;
  }

  @keyframes shine {
    0%,
    100% {
      background-position: left;
    }
    50% {
      background-position: right;
    }
  }

  section > div {
    width: 100%;
    white-space: nowrap;
    overflow: hidden;
    line-height: 2.25rem;
  }

  .number {
    font-size: var(--sl-font-size-x-small);
  }

  .name {
    text-overflow: ellipsis;
  }

  .score {
    justify-self: end;

    display: flex;
    align-items: center;
    gap: var(--sl-spacing-x-small);

    text-align: right;
    font-weight: var(--sl-font-weight-bold);
    font-size: var(--sl-font-size-medium);
  }
</style>
