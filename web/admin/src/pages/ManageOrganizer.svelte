<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import RelativeTime from "@/components/RelativeTime.svelte";
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/copy-button/copy-button.js";
  import {
    EmptyState,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import type { OrganizerInvite, User } from "@climblive/lib/models";
  import {
    createOrganizerInviteMutation,
    getOrganizerInvitesQuery,
    getOrganizerQuery,
    getSelfQuery,
    getUsersByOrganizerQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { isAfter } from "date-fns";
  import { navigate } from "svelte-routing";
  import CreateOrganizer from "./CreateOrganizer.svelte";
  import DeleteInvite from "./DeleteInvite.svelte";

  interface Props {
    organizerId: number;
  }

  const { organizerId }: Props = $props();

  const invitesQuery = $derived(getOrganizerInvitesQuery(organizerId));
  const createInvite = $derived(createOrganizerInviteMutation(organizerId));
  const usersQuery = $derived(getUsersByOrganizerQuery(organizerId));
  const selfQuery = $derived(getSelfQuery());
  const organizerQuery = $derived(getOrganizerQuery(organizerId));

  const invites = $derived(invitesQuery.data);
  const users = $derived(usersQuery.data);
  const self = $derived(selfQuery.data);
  const organizer = $derived(organizerQuery.data);

  const userColumns: ColumnDefinition<User>[] = [
    {
      label: "User",
      mobile: true,
      render: renderUsername,
      width: "1fr",
    },
  ];

  const inviteColumns: ColumnDefinition<OrganizerInvite>[] = [
    {
      label: "Invite link",
      mobile: true,
      render: renderCopyLink,
      width: "1fr",
    },
    {
      label: "Expires",
      mobile: true,
      render: renderExpiresAt,
      width: "max-content",
    },
    {
      mobile: true,
      render: renderControls,
      width: "max-content",
      align: "right",
    },
  ];

  const handleCreateInvite = () => {
    createInvite.mutate(undefined, {
      onError: () => toastError("Failed to create invite."),
    });
  };
</script>

{#snippet renderUsername({ id, username }: User)}
  {username}
  {#if id === self?.id}
    <wa-badge variant="brand" pill>Me</wa-badge>
  {/if}
{/snippet}

{#snippet renderCopyLink({ id }: OrganizerInvite)}
  {@const link = `${location.protocol}//${location.host}/admin/invites/${id}`}
  <wa-copy-button value={link}></wa-copy-button>
  {link}
{/snippet}

{#snippet renderExpiresAt({ expiresAt }: OrganizerInvite)}
  {@const expired = isAfter(new Date(), expiresAt)}

  <span class={{ expired }}>
    <RelativeTime time={expiresAt} />
  </span>
{/snippet}

{#snippet renderControls({ id }: OrganizerInvite)}
  <DeleteInvite inviteId={id}>
    {#snippet children({ deleteInvite })}
      <wa-button
        size="small"
        variant="danger"
        appearance="plain"
        onclick={deleteInvite}
      >
        <wa-icon name="trash" label={`Delete invite ${id}`}></wa-icon>
      </wa-button>
    {/snippet}
  </DeleteInvite>
{/snippet}

{#snippet createButton()}
  <wa-button
    variant="neutral"
    appearance="accent"
    onclick={handleCreateInvite}
    loading={createInvite.isPending}>Create invite</wa-button
  >
{/snippet}

<section>
  {#if invites === undefined || users === undefined || organizer === undefined}
    <Loader />
  {:else}
    <wa-breadcrumb>
      <wa-breadcrumb-item onclick={() => navigate("./")}
        ><wa-icon name="home"></wa-icon></wa-breadcrumb-item
      >
      <wa-breadcrumb-item
        onclick={() => navigate(`./organizers/${organizerId}`)}
        >{organizer.name}</wa-breadcrumb-item
      >
    </wa-breadcrumb>

    <CreateOrganizer>
      {#snippet children({ createOrganizer })}
        <wa-button
          size="small"
          variant="neutral"
          appearance="accent"
          onclick={createOrganizer}
        >
          <wa-icon slot="start" name="plus"></wa-icon>
          New organizer
        </wa-button>
      {/snippet}
    </CreateOrganizer>

    <h2>Co-organizers</h2>
    <Table columns={userColumns} data={users} getId={({ id }) => id}></Table>

    <h2>Invites</h2>
    {#if invites.length > 0}
      {@render createButton()}
    {/if}

    {#if invites.length > 0}
      <Table columns={inviteColumns} data={invites} getId={({ id }) => id}
      ></Table>
    {:else}
      <EmptyState
        title="No invites yet"
        description="Create invites to allow other users to join as co-organizers."
      >
        {#snippet actions()}
          {@render createButton()}
        {/snippet}
      </EmptyState>
    {/if}
  {/if}
</section>

<style>
  section {
    display: flex;
    flex-direction: column;
    align-items: start;
    gap: var(--wa-space-m);
  }

  .expired {
    color: var(--wa-color-danger);
  }
</style>
