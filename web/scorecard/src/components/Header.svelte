<script lang="ts">
  import { Score, Timer } from "@climblive/lib/components";
  import { type ContestState } from "@climblive/lib/types";
  import { ordinalSuperscript } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  interface Props {
    registrationCode: string;
    contestName: string;
    compClassName: string | undefined;
    contenderName: string | undefined;
    score: number;
    placement: number | undefined;
    contestState: ContestState;
    startTime: Date;
    endTime: Date;
  }

  let {
    registrationCode,
    contestName,
    compClassName,
    contenderName,
    score,
    placement,
    contestState,
    startTime,
    endTime,
  }: Props = $props();
</script>

<header>
  <wa-button
    size="small"
    onclick={() => navigate(`/${registrationCode}/edit`)}
    disabled={contestState === "ENDED"}
    appearance="plain"
  >
    <wa-icon name="gear" label="Edit"></wa-icon>
  </wa-button>
  <h1>{contestName}</h1>
  <p class="contender-name">
    {contenderName} <span class="contender-class">{compClassName}</span>
  </p>
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
    {#if contestState === "NOT_STARTED"}
      <Timer align="right" endTime={startTime} label="Time until start" />
    {:else}
      <Timer align="right" {endTime} label="Time remaining" />
    {/if}
  </div>
</header>

<style>
  header {
    background-color: var(--wa-color-brand-fill-normal);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-brand-border-normal);
    border-radius: var(--wa-border-radius-m);
    padding: var(--wa-space-s);
    color: var(--wa-color-brand-on-normal);
    position: relative;

    & wa-button {
      position: absolute;
      top: var(--wa-space-m);
      right: var(--wa-space-s);
      color: inherit;

      &::part(label) {
        color: var(--wa-color-brand-on-normal);
      }
    }

    & wa-button::part(base) {
      padding: 0;
    }

    & h1,
    & .contender-name {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    & h1 {
      margin: 0;
      font-size: var(--wa-font-size-l);
      width: calc(100% - 2rem);
      line-height: var(--wa-line-height-condensed);
    }

    & .contender-name {
      margin: 0;
      line-height: var(--wa-line-height-condensed);
    }

    & .contender-class {
      font-weight: var(--wa-font-weight-bold);
      font-size: var(--wa-font-size-xs);
    }

    & .score {
      & > span {
        font-weight: var(--wa-font-weight-bold);
        font-size: var(--wa-font-size-l);
      }

      & > :not(span) {
        font-size: var(--wa-font-size-xs);
        font-weight: var(--wa-font-weight-normal);
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
