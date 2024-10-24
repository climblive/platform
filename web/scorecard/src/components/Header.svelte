<script lang="ts">
  import type { ContestState } from "@/types/state";
  import { Score, Timer } from "@climblive/lib/components";
  import { asOrdinal } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  export let registrationCode: string;
  export let contestName: string;
  export let compClassName: string | undefined;
  export let contenderName: string | undefined;
  export let contenderClub: string | undefined;
  export let score: number;
  export let placement: number | undefined;
  export let state: ContestState;
  export let startTime: Date;
  export let endTime: Date;
  export let disabled: boolean;
</script>

<header>
  <sl-icon-button
    name="gear"
    label="Edit"
    on:click={() => navigate(`/${registrationCode}/edit`)}
    {disabled}
  >
  </sl-icon-button>
  <h1>{contestName}</h1>
  <span class="contender-name">{contenderName}</span>
  <span class="contender-club">{contenderClub}</span>
  <span class="contender-class">{compClassName}</span>
  <div class="lower">
    <div class="score">
      <Score value={score} />
      <span>{placement ? `${asOrdinal(placement)} place` : "-"}</span>
    </div>
    <div class="timer">
      {#if state === "NOT_STARTED"}
        <Timer endTime={startTime} />
        <span class="footer">Time until start</span>
      {:else}
        <Timer {endTime} />
        <span class="footer">Time remaining</span>
      {/if}
    </div>
  </div>
</header>

<style>
  header {
    background: linear-gradient(
      45deg,
      var(--sl-color-primary-500),
      var(--sl-color-primary-700)
    );
    border-radius: var(--sl-border-radius-small);
    padding: var(--sl-spacing-small);
    color: var(--sl-color-primary-100);
    position: relative;

    & sl-icon-button {
      position: absolute;
      top: var(--sl-spacing-medium);
      right: var(--sl-spacing-small);
      font-size: var(--sl-font-size-medium);
      color: inherit;
    }

    & sl-icon-button::part(base) {
      padding: 0;
    }

    & h1 {
      margin: 0;
      font-size: var(--sl-font-size-x-large);
    }

    & .contender-class {
      font-weight: var(--sl-font-weight-bold);
      font-size: var(--sl-font-size-x-small);
    }

    & .score {
      & > :not(span) {
        font-weight: var(--sl-font-weight-bold);
        font-size: var(--sl-font-size-2x-large);
      }

      & > span {
        font-size: var(--sl-font-size-x-small);
        font-weight: var(--sl-font-weight-normal);
      }
    }

    & .lower {
      margin-top: 1rem;
      display: flex;
      justify-content: space-between;
      align-items: end;
    }

    & .timer {
      font-weight: var(--sl-font-weight-bold);
      font-size: var(--sl-font-small);

      & > * {
        display: block;
      }

      & .footer {
        font-size: var(--sl-font-size-x-small);
        font-weight: var(--sl-font-weight-normal);
      }

      text-align: right;
    }
  }
</style>
