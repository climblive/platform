<script lang="ts">
  import { getContestsByOrganizerQuery } from "@climblive/lib/queries";
  import { Link } from "svelte-routing";

  interface Props {
    organizerId: number;
  }

  let { organizerId }: Props = $props();

  const contestsQuery = getContestsByOrganizerQuery(organizerId);

  let contests = $derived($contestsQuery.data);
</script>

<ul>
  {#if contests}
    {#each contests as contest (contest.id)}
      <li>
        <Link to="contests/{contest.id}">{contest.name}</Link>
      </li>
    {/each}
  {/if}
</ul>
