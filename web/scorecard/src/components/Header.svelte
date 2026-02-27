<script lang="ts">
  import { Logo, Score, Timer } from "@climblive/lib/components";
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
    disqualified: boolean;
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
    disqualified,
  }: Props = $props();
</script>

<header>
  <div class="title-row">
    <div class="logo"><Logo /></div>

    <h1>
      {contestName}
    </h1>

    <wa-button
      size="small"
      onclick={() => navigate(`/${registrationCode}/edit`)}
      disabled={contestState === "ENDED"}
      appearance="plain"
    >
      <wa-icon name="gear" label="Edit"></wa-icon>
    </wa-button>
  </div>
  <p class="subtitle-row">
    <span class="contender-name">{contenderName}</span> •
    <span class="contender-class">{compClassName}</span>
  </p>
  <div class="lower">
    <div class="score">
      <span>
        {#if disqualified}
          Disqualified
        {:else if placement}
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
    background-color: var(--wa-color-surface-raised);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-surface-border);
    border-radius: var(--wa-border-radius-m);
    padding: var(--wa-space-s);
    box-shadow: var(--wa-shadow-s);

    & .title-row {
      display: flex;
      align-items: center;
      gap: var(--wa-space-xs);

      & wa-button {
        color: inherit;
      }
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
      line-height: var(--wa-line-height-condensed);
    }

    & .subtitle-row {
      margin: 0;
      line-height: var(--wa-line-height-condensed);
      margin-block-start: calc(-1 * var(--wa-space-2xs));
    }

    & .contender-name {
      font-weight: var(--wa-font-weight-bold);
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

  .logo {
    width: 1.5rem;
    height: 1.5rem;
    flex-shrink: 0;
  }
</style>
