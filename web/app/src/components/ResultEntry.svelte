<script lang="ts">
  import type { ScoreboardContender } from "@climblive/shared/models";
  import { asOrdinal } from "@climblive/shared/utils";
  import Score from "./Score.svelte";

  export let placement: number;
  export let finalist: boolean;
  export let contender: ScoreboardContender;
</script>

<section data-finalist={finalist}>
  <span class="number">{asOrdinal(placement)}</span>
  <div>{contender.contenderName}</div>
  <div class="score"><Score value={contender.qualifyingScore} /></div>
</section>

<style>
  section {
    height: 2.25rem;
    background-color: var(--sl-color-primary-600);
    color: var(--sl-color-primary-200);
    display: grid;
    padding-inline: var(--sl-spacing-small);
    border-radius: var(--sl-border-radius-medium);
    grid-template-columns: 2rem 3fr 1fr;
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
  }

  .number {
    font-size: var(--sl-font-size-x-small);
  }

  .score {
    justify-self: end;
    text-align: right;
    font-weight: var(--sl-font-weight-bold);
    font-size: var(--sl-font-size-medium);
  }
</style>
