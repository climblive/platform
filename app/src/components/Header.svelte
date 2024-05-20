<script lang="ts">
  import { navigate } from "svelte-routing";
  import Score from "./Score.svelte";
  import Timer from "./Timer.svelte";

  export let registrationCode: string;
  export let contestName: string;
  export let compClassName: string | undefined;
  export let contenderName: string | undefined;
  export let contenderClub: string | undefined;
  export let points: number;
  export let endTime: Date;
  export let qualifyingProblems: number;
</script>

<header>
  <sl-icon-button
    name="gear"
    label="Edit"
    on:click={() => navigate(`/${registrationCode}/edit`)}
  >
  </sl-icon-button>
  <h1>{contestName}</h1>
  <span class="contender-name">{contenderName}</span>
  <span class="contender-club">{contenderClub}</span>
  <span class="contender-class">{compClassName}</span>
  <div class="lower">
    <div class="points">
      <Score value={points} size="large" />
      <span>{qualifyingProblems} hardest</span>
    </div>
    <div class="timer">
      <Timer {endTime} />
      <span>Time remaining</span>
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

    & .points {
      font-weight: var(--sl-font-weight-bold);

      & span {
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

      & span {
        font-size: var(--sl-font-size-x-small);
        font-weight: var(--sl-font-weight-normal);
      }

      text-align: right;
    }
  }
</style>
