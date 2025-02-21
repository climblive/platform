<script lang="ts">
  import { getContestsByOrganizerQuery } from "@climblive/lib/queries";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import { Link, navigate } from "svelte-routing";

  interface Props {
    organizerId: number;
  }

  let { organizerId }: Props = $props();

  const contestsQuery = getContestsByOrganizerQuery(organizerId);

  let contests = $derived($contestsQuery.data);
</script>

<sl-button
  size="large"
  variant="primary"
  onclick={() => navigate(`organizers/${organizerId}/contests/new`)}
  >Create</sl-button
>

<ul>
  {#if contests}
    {#each contests as contest (contest.id)}
      <li>
        <Link to="contests/{contest.id}">{contest.name}</Link>
      </li>
    {/each}
  {/if}
</ul>
