<script lang="ts">
  import { Score, Timer } from "@climblive/lib/components";
  import { type ContestState } from "@climblive/lib/types";
  import { ordinalSuperscript } from "@climblive/lib/utils";
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
</script>

<header>
  <sl-icon-button
    name="gear"
    label="Edit"
    on:click={() => navigate(`/${registrationCode}/edit`)}
    disabled={state === "ENDED"}
  >
  </sl-icon-button>
  <h1>{contestName}</h1>
  <p class="contender-name">
    {contenderName} <span class="contender-class">{compClassName}</span>
  </p>
  <p class="contender-club">{contenderClub}</p>
  <div class="lower">
    <div class="score">
      <span>
        {#if placement}
          {placement}<sup>{ordinalSuperscript(placement)}</sup>
        {:else}
          -
        {/if}
      </span>
      <Score value={score} />
    </div>
    {#if state === "NOT_STARTED"}
      <Timer endTime={startTime} label="Time until start" />
    {:else}
      <Timer {endTime} label="Time remaining" />
    {/if}
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

    & h1,
    & .contender-name,
    & .contender-club {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    & h1 {
      margin: 0;
      font-size: var(--sl-font-size-large);
      width: calc(100% - 2rem);
      line-height: var(--sl-line-height-dense);
    }

    & .contender-name,
    & .contender-club {
      margin: 0;
      line-height: var(--sl-line-height-dense);
    }

    & .contender-club {
      font-size: var(--sl-font-size-x-small);
    }

    & .contender-class {
      font-weight: var(--sl-font-weight-bold);
      font-size: var(--sl-font-size-x-small);
    }

    & .score {
      & > span {
        font-weight: var(--sl-font-weight-bold);
        font-size: var(--sl-font-size-large);
      }

      & > :not(span) {
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
  }
</style>
