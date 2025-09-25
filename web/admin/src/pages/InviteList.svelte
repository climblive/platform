<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import RelativeTime from "@/components/RelativeTime.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/copy-button/copy-button.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import type { OrganizerInvite } from "@climblive/lib/models";
  import {
    createOrganizerInviteMutation,
    getOrganizerInvitesQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";
  import DeleteInvite from "./DeleteInvite.svelte";

  interface Props {
    organizerId: number;
  }

  const { organizerId }: Props = $props();

  const invitesQuery = $derived(getOrganizerInvitesQuery(organizerId));
  const createInvite = $derived(createOrganizerInviteMutation(organizerId));

  let invites = $derived($invitesQuery.data);

  const columns: ColumnDefinition<OrganizerInvite>[] = [
    {
      label: "Invite",
      mobile: true,
      render: renderCopyLink,
      width: "1fr",
    },
    {
      label: "Expires at",
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
    $createInvite.mutate(undefined, {
      onError: () => toastError("Failed to create invite."),
    });
  };
</script>

{#snippet renderCopyLink({ id }: OrganizerInvite)}
  <wa-copy-button
    value={`${location.protocol}//${location.host}/admin/invites/${id}`}
  ></wa-copy-button>
{/snippet}

{#snippet renderExpiresAt({ expiresAt }: OrganizerInvite)}
  <RelativeTime time={expiresAt} />
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

<wa-button
  appearance="plain"
  onclick={() => navigate(`./organizers/${organizerId}`)}
  >Back<wa-icon name="arrow-left" slot="start"></wa-icon></wa-button
>

<section>
  <wa-button
    variant="neutral"
    appearance="accent"
    onclick={handleCreateInvite}
    loading={$createInvite.isPending}>Create invite</wa-button
  >

  {#if invites === undefined}
    <Loader />
  {:else if invites.length > 0}
    <Table {columns} data={invites} getId={({ id }) => id}></Table>
  {/if}
</section>

<style>
  section {
    display: flex;
    flex-direction: column;
    align-items: start;
    gap: var(--wa-space-m);
  }
</style>
