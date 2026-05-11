<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
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
  import { getContext } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Authenticator } from "../authenticator.svelte";

  const selfQuery = $derived(getSelfQuery());
  const self = $derived(selfQuery.data);
  const authenticator = getContext<Authenticator>("authenticator");

  const organizerColumns: ColumnDefinition<Organizer>[] = [
    {
      label: "Name",
      mobile: true,
      render: renderOrganizerName,
      width: "1fr",
    },
    {
      label: "",
      mobile: true,
      render: renderOrganizerActions,
      width: "max-content",
      align: "right",
    },
  ];
</script>

{#snippet renderOrganizerName({ name }: Organizer)}
  {name}
{/snippet}

{#snippet renderOrganizerActions({ id, name }: Organizer)}
  <wa-button
    size="s"
    appearance="plain"
    aria-label={`Open settings for organizer ${name}`}
    onclick={() => navigate(`/admin/organizers/${id}`)}
  >
    <wa-icon slot="start" name="gear"></wa-icon>
    Settings
  </wa-button>
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
      </div>

      <div>
        <h2>Username</h2>
        <div class="username-row">
          <p>{self.username}</p>
          <wa-button
            variant="neutral"
            appearance="outlined"
            size="s"
            onclick={authenticator.logout}
          >
            <wa-icon slot="start" name="right-from-bracket"></wa-icon>
            Sign out
          </wa-button>
        </div>
      </div>

      <div>
        <h2>Organizers</h2>
        {#if self.organizers.length > 0}
          <div class="organizers-table">
            <Table
              columns={organizerColumns}
              data={self.organizers}
              getId={({ id }) => id}
            ></Table>
          </div>
        {:else}
          <EmptyState
            title="No organizers yet"
            description="You are not part of any organizers."
          />
        {/if}
      </div>

      <div class="remove-account">
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

  p {
    margin: 0;
  }

  .username-row {
    display: flex;
    align-items: center;
    gap: var(--wa-space-s);
  }

  .remove-account wa-callout {
    margin-block-start: var(--wa-space-s);
  }
</style>
