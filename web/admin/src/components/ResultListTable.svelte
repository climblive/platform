<script lang="ts">
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import type WaInput from "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import {
    CompetitorName,
    Score,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import { value } from "@climblive/lib/forms";
  import type { ScoreboardEntry } from "@climblive/lib/models";
  import { getCategoriesQuery } from "@climblive/lib/queries";
  import { formatAgeSpan, ordinalSuperscript } from "@climblive/lib/utils";
  import { Link } from "svelte-routing";
  import type { Readable } from "svelte/store";
  import Loader from "./Loader.svelte";

  interface Props {
    competitionId: number;
    scoreboard: Readable<Map<number, ScoreboardEntry[]>>;
    loading: boolean;
  }

  const { competitionId, scoreboard, loading }: Props = $props();

  let tableData = $state<ScoreboardEntry[]>([]);

  let categorySelector: WaSelect | undefined = $state();
  let quickFilter: WaInput | undefined = $state();

  const categoriesQuery = $derived(getCategoriesQuery(competitionId));

  const categories = $derived(categoriesQuery.data);

  let filterText = $state<string>();
  let selectedCategoryId: number | undefined = $state();

  const competitorCounts = $derived(
    new Map(
      Array.from($scoreboard.entries()).map(([categoryId, entries]) => [
        categoryId,
        entries.length,
      ]),
    ),
  );

  $effect(() => {
    if (
      categories &&
      categories.length > 0 &&
      selectedCategoryId === undefined
    ) {
      selectedCategoryId = categories[0].id;
    }
  });

  $effect(() => {
    if (selectedCategoryId === undefined) {
      return;
    }

    let results = [...($scoreboard.get(selectedCategoryId) ?? [])];

    if (filterText) {
      const search = filterText.toLowerCase();

      results = results.filter(({ name }) =>
        name.toLowerCase().includes(search),
      );
    }

    results.sort(
      (a: ScoreboardEntry, b: ScoreboardEntry) =>
        (a.result?.rankOrder ?? Infinity) - (b.result?.rankOrder ?? Infinity),
    );

    tableData = results;
  });

  const columns: ColumnDefinition<ScoreboardEntry>[] = [
    {
      label: "Name",
      mobile: true,
      render: renderName,
      width: "3fr",
    },
    {
      label: "Score",
      mobile: false,
      render: renderScore,
      width: "max-content",
      align: "right",
    },
    {
      label: "Placement",
      mobile: true,
      render: renderPlacement,
      width: "max-content",
      align: "right",
    },
    {
      label: "Finalist",
      mobile: true,
      render: renderFinalist,
      width: "max-content",
      align: "left",
    },
  ];
</script>

{#snippet renderName({
  competitorId,
  name,
  disqualified,
  scrubbedAt,
}: ScoreboardEntry)}
  <Link to={`./competitors/${competitorId}`}>
    {#if disqualified}
      <del><CompetitorName id={competitorId} {name} {scrubbedAt} /></del>
    {:else}
      <CompetitorName id={competitorId} {name} {scrubbedAt} />
    {/if}
  </Link>
{/snippet}

{#snippet renderScore({ result }: ScoreboardEntry)}
  {#if result}
    <Score value={result.score}></Score>
  {/if}
{/snippet}

{#snippet renderPlacement({ result }: ScoreboardEntry)}
  {#if result?.placement}
    {result.placement}<sup>{ordinalSuperscript(result.placement)}</sup>
  {:else}
    <wa-icon name="minus"></wa-icon>
  {/if}
{/snippet}

{#snippet renderFinalist({ result }: ScoreboardEntry)}
  <wa-icon name={result?.finalist ? "medal" : "minus"}></wa-icon>
{/snippet}

<h2>Results</h2>

<section>
  {#if categories}
    <div class="controls">
      <wa-input
        bind:this={quickFilter}
        size="s"
        label="Quick filter"
        placeholder="Search by name..."
        oninput={() => {
          filterText = quickFilter?.value ?? "";
        }}
        with-clear
      ></wa-input>
      <wa-select
        bind:this={categorySelector}
        size="s"
        label="Category filter"
        {@attach value(selectedCategoryId)}
        onchange={() => {
          selectedCategoryId = Number(categorySelector?.value);
        }}
      >
        {#each categories as category (category.id)}
          {@const count = competitorCounts.get(category.id)}

          <wa-option value={category.id} label={category.name}>
            <div class="label">
              {category.name}
              {#if count}
                <wa-badge pill variant="neutral">{count}</wa-badge>
              {/if}
            </div>
            <small>Ages {formatAgeSpan(category.minAge, category.maxAge)}</small
            >
          </wa-option>
        {/each}
      </wa-select>
    </div>
  {/if}

  {#if loading}
    <Loader />
  {:else}
    <Table {columns} data={tableData} getId={({ competitorId }) => competitorId}
    ></Table>
  {/if}
</section>

<style>
  .label {
    display: flex;
    justify-content: start;
    align-items: center;
    gap: var(--wa-space-xs);
    flex-grow: 0;
  }

  .controls {
    display: flex;
    gap: var(--wa-space-m);
    justify-content: space-evenly;

    & > * {
      width: 100%;
    }
  }

  section {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }
</style>
