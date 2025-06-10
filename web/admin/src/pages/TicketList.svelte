<script lang="ts">
  import {
    getContendersByContestQuery,
    getContestQuery,
  } from "@climblive/lib/queries";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import "@shoelace-style/shoelace/dist/components/details/details.js";
  import "@shoelace-style/shoelace/dist/components/icon/icon.js";
  import "@shoelace-style/shoelace/dist/components/tab-group/tab-group.js";
  import "@shoelace-style/shoelace/dist/components/tab-panel/tab-panel.js";
  import "@shoelace-style/shoelace/dist/components/tab/tab.js";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const contestQuery = $derived(getContestQuery(contestId));
  const contendersQuery = $derived(getContendersByContestQuery(contestId));

  const contest = $derived($contestQuery.data);
  const contenders = $derived($contendersQuery.data);
</script>

<main>
  {#if contest && contenders}
    {#each contenders as contender (contender.id)}
      <section>
        <sl-qr-code
          size="64"
          value={`${location.protocol}//${location.host}/${contender.registrationCode}`}
        ></sl-qr-code>
        {contender.registrationCode}
      </section>
    {/each}
  {/if}
</main>

<style>
  @page {
    size: a4 portrait;
    margin: 2cm;
  }

  section {
    border: 1px solid var(--sl-color-primary-600);
    border-radius: var(--sl-border-radius-medium);
    margin-block: var(--sl-spacing-medium);
    padding: var(--sl-spacing-medium);
  }

  section:nth-child(7n) {
    break-after: page;
  }
</style>
