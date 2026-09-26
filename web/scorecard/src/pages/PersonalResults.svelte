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
  <nav aria-label="Personal results">
    <wa-button
      appearance="plain"
      size="s"
      onclick={() => navigate(`/${registrationCode}`)}
    >
      <wa-icon slot="start" name="arrow-left"></wa-icon>
      Back to scorecard
    </wa-button>
  </nav>

  <div class="poster-container">
    <article class="poster" aria-label="Personal results">
      <header>
        <div class="branding">
          <div class="logo" role="img" aria-label="ClimbLive">
            <FullLogo />
          </div>
          <span class="eyebrow">Personal results</span>
        </div>
        <p class="competition">{contest.name}</p>
        <h1>
          <ContenderName
            id={contender.id}
            name={contender.name}
            scrubbedAt={contender.scrubbedAt}
          />
        </h1>
        <p class="category">{compClassName}</p>
      </header>

      <dl class="stats">
        <div>
          <dt>Score</dt>
          <dd class:long-score={score.length > 8}>{score}</dd>
        </div>
        <div>
          <dt>Placement</dt>
          <dd>
            {#if contender.disqualified}
              <span class="unranked">Disqualified</span>
            {:else if placement && placement > 0}
              {placement}<sup>{ordinalSuperscript(placement)}</sup>
            {:else}
              <span class="unranked">Unranked</span>
            {/if}
          </dd>
        </div>
      </dl>

      <section aria-labelledby="hardest-problems">
        <div class="section-heading">
          <h2 id="hardest-problems">Top 5 hardest problems</h2>
          <span class="eyebrow">Completed</span>
        </div>
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
                <span class="problem-number">Problem {problem.number}</span>
                <span class="ascent">
                  {#if problem.tick.attemptsTop === 1}
                    <wa-icon name="bolt"></wa-icon> Flash
                  {:else}
                    Top
                  {/if}
                </span>
                {#if contest.usePoints}
                  <strong class="points">
                    {problem.pointValue
                      ? `${problem.pointValue.current}p`
                      : "—"}
                  </strong>
                {/if}
              </li>
            {/each}
          </ol>
        {:else}
          <p class="empty">
            Your hardest tops will appear here once you complete a problem.
          </p>
        {/if}
      </section>

      <footer>
        <span
          >{contest.usePoints
            ? "Ranked by top point value"
            : "Ranked by problem point value"}</span
        >
        <strong>climblive.app</strong>
      </footer>
    </article>
  </div>
  <p class="hint">Take a screenshot of your results to share on Instagram.</p>
</main>

<style>
  main {
    padding: var(--wa-space-m);
  }

  nav {
    margin-block-end: var(--wa-space-m);
  }

  .poster-container {
    container-type: inline-size;
  }

  .poster {
    aspect-ratio: 4 / 5;
    display: flex;
    flex-direction: column;
    gap: 3cqi;
    padding: 5cqi;
    border-radius: var(--wa-border-radius-l);
    border: var(--wa-border-width-s) solid var(--wa-color-brand-80);
    background: linear-gradient(
      145deg,
      var(--wa-color-neutral-95),
      var(--wa-color-brand-95)
    );
    color: var(--wa-color-neutral-10);
    font-size: clamp(var(--wa-font-size-2xs), 3cqi, var(--wa-font-size-m));
    overflow-wrap: anywhere;
  }

  .branding,
  .section-heading,
  footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--wa-space-s);
  }

  .branding {
    margin-block-end: 4cqi;
  }

  .logo {
    height: 5cqi;
    flex-shrink: 0;
  }

  .eyebrow,
  dt {
    text-transform: uppercase;
    font-size: 0.8em;
    font-weight: var(--wa-font-weight-semibold);
    letter-spacing: 0.1em;
  }

  .eyebrow,
  .category,
  .ascent,
  footer {
    color: var(--wa-color-neutral-40);
  }

  p,
  h1,
  h2,
  dl,
  dd {
    margin: 0;
  }

  .competition {
    font-weight: var(--wa-font-weight-semibold);
    color: var(--wa-color-brand-40);
  }

  h1 {
    margin-block: var(--wa-space-xs);
    font-size: 7cqi;
    line-height: var(--wa-line-height-condensed);
    font-weight: var(--wa-font-weight-bold);
  }

  .stats {
    display: grid;
    grid-template-columns: minmax(0, 1.25fr) minmax(0, 1fr);
    gap: var(--wa-space-m);
    padding-block: 3cqi;
    border-block: var(--wa-border-width-s) solid var(--wa-color-brand-80);
  }

  dt {
    margin-block-end: var(--wa-space-xs);
    color: var(--wa-color-brand-40);
  }

  dd {
    font-size: 11cqi;
    line-height: var(--wa-line-height-condensed);
    font-weight: var(--wa-font-weight-bold);
    font-variant-numeric: tabular-nums;
  }

  dd.long-score {
    font-size: 6cqi;
  }

  sup {
    font-size: 0.4em;
  }

  .unranked {
    display: block;
    font-size: 4cqi;
  }

  h2 {
    font-size: 1em;
    font-weight: var(--wa-font-weight-bold);
  }

  ol {
    list-style: none;
    padding: 0;
    margin: var(--wa-space-s) 0 0;
  }

  li {
    display: flex;
    align-items: center;
    gap: var(--wa-space-xs);
    padding-block: 1.4cqi;
    border-bottom: var(--wa-border-width-s) solid var(--wa-color-brand-90);
  }

  .problem-number {
    flex: 1;
    font-weight: var(--wa-font-weight-semibold);
  }

  .ascent {
    display: flex;
    align-items: center;
    gap: var(--wa-space-2xs);
    font-size: 0.85em;
  }

  .ascent wa-icon {
    color: var(--wa-color-brand-40);
  }

  .points {
    min-width: 4em;
    text-align: right;
    font-variant-numeric: tabular-nums;
  }

  .empty {
    padding-block: var(--wa-space-l);
    line-height: var(--wa-line-height-normal);
  }

  footer {
    margin-block-start: auto;
    padding-block-start: var(--wa-space-s);
    font-size: 0.75em;
  }

  .hint {
    margin-block-start: var(--wa-space-m);
    color: var(--wa-color-text-quiet);
    font-size: var(--wa-font-size-s);
    text-align: center;
  }
</style>
