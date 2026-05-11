<script lang="ts">
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import Loader from "@/components/Loader.svelte";
  import {
    EmptyState,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import type { Organizer } from "@climblive/lib/models";
  import { getSelfQuery } from "@climblive/lib/queries";
  import { navigate } from "svelte-routing";

  const selfQuery = $derived(getSelfQuery());
  const self = $derived(selfQuery.data);

  const organizerColumns: ColumnDefinition<Organizer>[] = [
    {
      label: "Organizer",
      mobile: true,
      render: renderOrganizerName,
      width: "1fr",
    },
  ];
</script>

{#snippet renderOrganizerName({ name }: Organizer)}
  {name}
{/snippet}

{#if self === undefined}
  <Loader />
{:else}
  <section>
    <wa-breadcrumb>
      <wa-breadcrumb-item onclick={() => navigate("./")}
        ><wa-icon name="home"></wa-icon></wa-breadcrumb-item
      >
      <wa-breadcrumb-item>Profile</wa-breadcrumb-item>
    </wa-breadcrumb>

    <div class="content">
      <div>
        <h1>Profile</h1>
        <p class="description">
          Manage your account details and support options.
        </p>
      </div>

      <div>
        <h2>Username</h2>
        <p>{self.username}</p>
      </div>

      <div>
        <h2>Organizers</h2>
        {#if self.organizers.length > 0}
          <Table
            columns={organizerColumns}
            data={self.organizers}
            getId={({ id }) => id}
          ></Table>
        {:else}
          <EmptyState
            title="No organizers yet"
            description="You are not part of any organizers."
          />
        {/if}
      </div>

      <div>
        <h2>Remove account</h2>
        <wa-callout variant="warning">
          <wa-icon slot="icon" name="triangle-exclamation"></wa-icon>
          To remove your account, please contact support at
          <a href="mailto:info@climblive.com">info@climblive.com</a>.
        </wa-callout>
      </div>
    </div>
  </section>
{/if}

<style>
  section {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }

  .content {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-xl);
  }

  .description {
    margin: 0;
    color: var(--wa-color-text-quiet);
  }

  h1,
  h2 {
    margin-block: 0;
  }

  p {
    margin-block: var(--wa-space-xs) 0;
  }
</style>
