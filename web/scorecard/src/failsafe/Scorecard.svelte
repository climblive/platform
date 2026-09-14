<script lang="ts">
  import { getContenderQuery, getContestQuery } from "@climblive/lib/queries";
  import EditProfile from "./EditProfile.svelte";
  import ProblemList from "./ProblemList.svelte";

  type Props = {
    contestId: number;
    contenderId: number;
  };

  const { contestId, contenderId }: Props = $props();

  const contenderQuery = $derived(getContenderQuery(contenderId));
  const contestQuery = $derived(getContestQuery(contestId));

  const contender = $derived(contenderQuery.data);
  const contest = $derived(contestQuery.data);
</script>

<h2>Profile</h2>
<EditProfile {contestId} {contenderId} />

{#if contender?.entered}
  <h2>Scorecard</h2>

  {#if contest?.usePoints === false}
    <div class="info">
      <strong>Important notice</strong>
      <span>
        In this basic version of the app, any result you log that is not a flash
        will be recorded as 999 attempts.
      </span>
    </div>
  {/if}

  <ProblemList {contestId} {contenderId}></ProblemList>
{/if}

<style>
  .info {
    padding: var(--wa-space-s);
    border: 1px solid var(--wa-color-warning-border-quiet);
    border-radius: 0.25rem;
    margin: 0;
    background-color: var(--wa-color-warning-fill-quiet);

    & strong {
      display: block;
      margin-block-end: var(--wa-space-2xs);
    }
  }
</style>
