<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/qr-code/qr-code.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import type { Contender } from "@climblive/lib/models";
  import { getContendersByContestQuery } from "@climblive/lib/queries";
  import { getApiUrl, ordinalSuperscript } from "@climblive/lib/utils";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const contendersQuery = $derived(getContendersByContestQuery(contestId));

  let contenders = $derived($contendersQuery.data);

  const columns: ColumnDefinition<Contender>[] = [
    {
      label: "Code",
      mobile: false,
      render: renderRegistrationCode,
      width: "max-content",
    },
    {
      label: "Name",
      mobile: true,
      render: renderName,
      width: "3fr",
    },
    {
      label: "Score",
      mobile: true,
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
  ];
</script>

{#snippet renderRegistrationCode({ registrationCode }: Contender)}
  <a href={`/${registrationCode}`} target="blank">
    {registrationCode}
  </a>
{/snippet}

{#snippet renderName({ name }: Contender)}
  {name}
{/snippet}

{#snippet renderScore({ score }: Contender)}
  {#if score}
    {score.score}
  {/if}
{/snippet}

{#snippet renderPlacement({ score }: Contender)}
  {#if score}
    {score.placement}<sup>{ordinalSuperscript(score.placement)}</sup>
  {/if}
{/snippet}

{#if contenders?.length}
  <a href={`${getApiUrl()}/contests/${contestId}/results`}>
    <wa-button appearance="outlined"
      >Download results
      <wa-icon name="download" slot="start"></wa-icon>
    </wa-button>
  </a>

  <Table {columns} data={contenders} getId={({ id }) => id}></Table>
{/if}
