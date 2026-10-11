<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import {
    ContenderName,
    FullLogo,
    HoldColorIndicator,
  } from "@climblive/lib/components";
  import type {
    Contender,
    Contest,
    PointValue,
    Problem,
    Tick,
  } from "@climblive/lib/models";
  import { ordinalSuperscript } from "@climblive/lib/utils";
  import { format, isSameDay } from "date-fns";
  import { sv } from "date-fns/locale";
  import { navigate } from "svelte-routing";

  interface Props {
    registrationCode: string;
    contest: Contest;
    contender: Contender;
    compClassName: string;
    score: string;
    placement: number | undefined;
    problems: Problem[];
    ticks: Tick[];
    pointValues: PointValue[] | undefined;
  }

  const {
    registrationCode,
    contest,
    contender,
    compClassName,
    score,
    placement,
    problems,
    ticks,
    pointValues,
  }: Props = $props();

  const hardestProblems = $derived.by(() => {
    const tops = new Map(
      ticks.filter((tick) => tick.top).map((tick) => [tick.problemId, tick]),
    );
    const values = new Map(
      (pointValues ?? []).map((value) => [value.problemId, value]),
    );

    return problems
      .flatMap((problem) => {
        const tick = tops.get(problem.id);

        return tick
          ? [{ ...problem, tick, pointValue: values.get(problem.id) }]
          : [];
      })
      .sort(
        (a, b) =>
          (b.pointValue?.top ?? b.pointsTop) -
            (a.pointValue?.top ?? a.pointsTop) || a.number - b.number,
      )
      .slice(0, 5);
  });
</script>

<svelte:head>
  <title>Personal results · ClimbLive</title>
</svelte:head>

<main>
  <wa-button
    appearance="plain"
    size="s"
    onclick={() => navigate(`/${registrationCode}`)}
  >
    <wa-icon slot="start" name="arrow-left"></wa-icon>
    Back to scorecard
  </wa-button>

  <article aria-label="Personal results">
    <header>
      <div class="logo" role="img" aria-label="ClimbLive">
        <FullLogo />
      </div>
      <h1>
        <ContenderName
          id={contender.id}
          name={contender.name}
          scrubbedAt={contender.scrubbedAt}
        />
      </h1>
      <p>{contest.name} / {compClassName}</p>
      {#if contest.timeBegin}
        <p class="date">
          <time datetime={contest.timeBegin.toISOString()}>
            {format(contest.timeBegin, "PPp", { locale: sv })}
          </time>
          {#if contest.timeEnd}
            –
            <time datetime={contest.timeEnd.toISOString()}>
              {format(
                contest.timeEnd,
                isSameDay(contest.timeBegin, contest.timeEnd) ? "p" : "PPp",
                { locale: sv },
              )}
            </time>
          {/if}
        </p>
      {/if}
    </header>

    <dl>
      <div>
        <dt>Score</dt>
        <dd class:long-score={score.length > 8}>{score}</dd>
      </div>
      <div>
        <dt>Placement</dt>
        <dd
          class:unranked={contender.disqualified || !placement || placement < 0}
        >
          {#if contender.disqualified}
            Disqualified
          {:else if placement && placement > 0}
            {placement}<sup>{ordinalSuperscript(placement)}</sup>
          {:else}
            Unranked
          {/if}
        </dd>
      </div>
    </dl>

    <section aria-labelledby="hardest-problems">
      <h2 id="hardest-problems">Top 5 hardest problems</h2>
      {#if hardestProblems.length > 0}
        <ol>
          {#each hardestProblems as problem (problem.id)}
            <li>
              <HoldColorIndicator
                primary={problem.holdColorPrimary}
                secondary={problem.holdColorSecondary}
                --height="var(--wa-font-size-m)"
                --width="var(--wa-font-size-m)"
              />
              <span class="problem-number">Problem #{problem.number}</span>
              <span class="ascent">
                {#if problem.tick.attemptsTop === 1}
                  <wa-icon name="bolt"></wa-icon> Flash
                {:else}
                  Top
                {/if}
              </span>
              {#if contest.usePoints}
                <strong>
                  {problem.pointValue ? `${problem.pointValue.current}p` : "—"}
                </strong>
              {/if}
            </li>
          {/each}
        </ol>
      {:else}
        <p>Your hardest tops will appear here once you complete a problem.</p>
      {/if}
    </section>

    <footer>climblive.app</footer>
  </article>
</main>

<style>
  main {
    container-type: inline-size;
    padding: var(--wa-space-m);
  }

  article {
    aspect-ratio: 4 / 5;
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-l);
    margin-block-start: var(--wa-space-m);
    padding: var(--wa-space-l);
    background-color: var(--wa-color-surface-raised);
    border: var(--wa-border-width-s) solid var(--wa-color-surface-border);
    border-radius: var(--wa-border-radius-m);
    overflow-wrap: anywhere;
  }

  .logo {
    height: var(--wa-font-size-l);
    margin-block-end: var(--wa-space-l);
  }

  h1,
  h2,
  p,
  dl,
  dd {
    margin: 0;
  }

  h1 {
    margin-block-end: var(--wa-space-2xs);
    font-size: clamp(var(--wa-font-size-xl), 7cqi, var(--wa-font-size-3xl));
  }

  header p,
  dt,
  .ascent,
  footer {
    color: var(--wa-color-text-quiet);
  }

  .date {
    margin-block-start: var(--wa-space-xs);
    font-size: var(--wa-font-size-s);
  }

  dl {
    display: grid;
    grid-template-columns: minmax(0, 1.25fr) minmax(0, 1fr);
    gap: var(--wa-space-m);
    padding-block: var(--wa-space-m);
    border-block: var(--wa-border-width-s) solid var(--wa-color-surface-border);
    font-variant-numeric: tabular-nums;
  }

  dt {
    font-size: var(--wa-font-size-xs);
  }

  dd {
    font-size: clamp(var(--wa-font-size-2xl), 11cqi, var(--wa-font-size-4xl));
    font-weight: var(--wa-font-weight-bold);
    line-height: 1;
  }

  dd.long-score {
    font-size: clamp(var(--wa-font-size-l), 6cqi, var(--wa-font-size-2xl));
  }

  dd.unranked {
    font-size: var(--wa-font-size-m);
  }

  sup {
    font-size: 0.4em;
  }

  h2 {
    margin-block-end: var(--wa-space-s);
    font-size: var(--wa-font-size-m);
  }

  ol {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  li {
    display: flex;
    align-items: center;
    gap: var(--wa-space-xs);
    padding-block: var(--wa-space-xs);
    border-bottom: var(--wa-border-width-s) solid var(--wa-color-surface-border);
    font-size: var(--wa-font-size-s);

    strong {
      min-width: 4em;
      text-align: right;
      font-variant-numeric: tabular-nums;
    }
  }

  .problem-number {
    flex: 1;
    font-weight: var(--wa-font-weight-semibold);
  }

  .ascent {
    font-size: var(--wa-font-size-xs);
  }

  footer {
    margin-block-start: auto;
    font-size: var(--wa-font-size-xs);
    font-weight: var(--wa-font-weight-semibold);
    text-align: right;
  }
</style>
