<script lang="ts">
  import { FullLogo } from "@climblive/lib/components";
  import type { Category, Problem } from "@climblive/lib/models";
  import { onMount, tick } from "svelte";

  interface Props {
    competitionName: string;
    categories: Category[];
    problems: Problem[];
  }

  const { competitionName, categories, problems }: Props = $props();

  let rowsPerColumn = $state(24);
  let tableContainers = $state<HTMLDivElement[]>([]);

  const sortedProblems = $derived(
    problems.toSorted((a, b) => a.number - b.number),
  );

  const pages = $derived.by(() => {
    const capacity = rowsPerColumn * 3;
    const pages: Problem[][][] = [];

    for (let start = 0; start < sortedProblems.length; start += capacity) {
      const page = sortedProblems.slice(start, start + capacity);
      const columnLength = Math.ceil(page.length / 3);
      const columns: Problem[][] = [];

      for (let column = 0; column < 3; column++) {
        columns.push(
          page.slice(column * columnLength, (column + 1) * columnLength),
        );
      }

      pages.push(columns);
    }

    return pages;
  });

  const updateRowsPerColumn = () => {
    const container = tableContainers[0];
    const heading = container?.querySelector("thead");
    const row = container?.querySelector("tbody tr");

    if (!container || !heading || !row) {
      return;
    }

    const availableHeight =
      container.clientHeight - heading.getBoundingClientRect().height - 1;
    rowsPerColumn = Math.max(
      1,
      Math.floor(availableHeight / row.getBoundingClientRect().height),
    );
  };

  export const preparePrint = async () => {
    updateRowsPerColumn();
    await tick();
  };

  onMount(() => {
    const observer = new ResizeObserver(updateRowsPerColumn);
    observer.observe(tableContainers[0]);
    updateRowsPerColumn();

    return () => observer.disconnect();
  });
</script>

{#each pages as columns, pageIndex (pageIndex)}
  <section class="sheet" aria-label={`Scorecard page ${pageIndex + 1}`}>
    <header>
      <div class="brand">
        <div class="logo" role="img" aria-label="ClimbLive"><FullLogo /></div>
        <span>Paper scorecard</span>
      </div>
      <h1>{competitionName}</h1>
      <div class="name">Name <span class="writing-line"></span></div>
      <div class="categories">
        <strong>Category</strong>
        {#each categories as category (category.id)}
          <span class="category">
            <span class="checkbox" aria-hidden="true"></span>{category.name}
          </span>
        {:else}
          <span class="writing-line"></span>
        {/each}
      </div>
      <p>
        Tick only your highest achievement for each problem. Flash = top on the
        first attempt. Leave unclimbed problems blank.
      </p>
    </header>

    <div class="problem-tables" bind:this={tableContainers[pageIndex]}>
      {#each columns as column, columnIndex (columnIndex)}
        {#if column.length > 0}
          {@const zone1 = column.some((problem) => problem.zone1Enabled)}
          {@const zone2 = column.some((problem) => problem.zone2Enabled)}
          <table
            aria-label={`Problems ${column[0].number}–${column[column.length - 1].number}`}
          >
            <thead>
              <tr>
                <th scope="col">No.</th>
                {#if zone1}<th scope="col">Z1</th>{/if}
                {#if zone2}<th scope="col">Z2</th>{/if}
                <th scope="col">Top</th>
                <th scope="col">Flash</th>
              </tr>
            </thead>
            <tbody>
              {#each column as problem (problem.id)}
                <tr>
                  <th scope="row">{problem.number}</th>
                  {#if zone1}
                    <td
                      aria-label={problem.zone1Enabled
                        ? "Zone 1"
                        : "Not available"}
                    >
                      {#if problem.zone1Enabled}
                        <span class="checkbox" aria-hidden="true"></span>
                      {:else}—{/if}
                    </td>
                  {/if}
                  {#if zone2}
                    <td
                      aria-label={problem.zone2Enabled
                        ? "Zone 2"
                        : "Not available"}
                    >
                      {#if problem.zone2Enabled}
                        <span class="checkbox" aria-hidden="true"></span>
                      {:else}—{/if}
                    </td>
                  {/if}
                  <td aria-label="Top"
                    ><span class="checkbox" aria-hidden="true"></span></td
                  >
                  <td aria-label="Flash"
                    ><span class="checkbox" aria-hidden="true"></span></td
                  >
                </tr>
              {/each}
            </tbody>
          </table>
        {/if}
      {/each}
    </div>

    <footer>
      <span>Return your scorecard to the organizer.</span>
      <span>Page {pageIndex + 1} of {pages.length}</span>
    </footer>
  </section>
{/each}

<style>
  .sheet {
    width: 190mm;
    height: 276mm;
    margin: 0 auto var(--wa-space-l);
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 4mm;
    background: white;
    color: black;
    font-family: Arial, sans-serif;
    font-size: 10pt;
    line-height: 1.25;
  }

  header {
    display: flex;
    flex-direction: column;
    gap: 4mm;
  }

  .brand {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 4mm;
  }

  .logo {
    height: 8mm;
    flex-shrink: 0;
    filter: grayscale(1);
  }

  h1 {
    margin: 0;
    font-size: 16pt;
    overflow-wrap: anywhere;
  }

  .name {
    display: flex;
    gap: 3mm;
    font-weight: bold;
  }

  .writing-line {
    display: inline-block;
    flex: 1;
    min-width: 30mm;
    min-height: 6mm;
    border-bottom: 0.25mm solid black;
  }

  .categories {
    display: flex;
    align-items: baseline;
    flex-wrap: wrap;
    gap: 3mm 5mm;
  }

  .category {
    display: inline-flex;
    align-items: center;
    gap: 2mm;
    max-width: 100%;
    overflow-wrap: anywhere;
  }

  .checkbox {
    display: inline-block;
    width: 3.5mm;
    height: 3.5mm;
    flex-shrink: 0;
    border: 0.25mm solid black;
    vertical-align: middle;
  }

  p {
    margin: 0;
    font-size: 9pt;
  }

  .problem-tables {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    align-items: start;
    gap: 4mm;
    flex: 1;
    min-height: 0;
  }

  table {
    width: 100%;
    table-layout: fixed;
    border-collapse: collapse;
    font-size: 9pt;
  }

  th,
  td {
    height: 8mm;
    padding: 0 1mm;
    border: 0.25mm solid black;
    text-align: center;
    vertical-align: middle;
  }

  footer {
    display: flex;
    justify-content: space-between;
    gap: 4mm;
    font-size: 8pt;
  }

  @media print {
    .sheet {
      margin: 0;
      break-after: page;

      &:last-child {
        break-after: auto;
      }
    }
  }
</style>
