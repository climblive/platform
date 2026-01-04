<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import { toastError, toastSuccess } from "@climblive/lib/utils";
  import type { UnlockRequest } from "@climblive/lib/models";
  import {
    getContestQuery,
    getOrganizerQuery,
    getPendingUnlockRequestsQuery,
    reviewUnlockRequestMutation,
  } from "@climblive/lib/queries";
  import { Link } from "svelte-routing";

  const pendingRequestsQuery = getPendingUnlockRequestsQuery();

  const handleApprove = async (request: UnlockRequest) => {
    const mutation = reviewUnlockRequestMutation(request.id);
    try {
      await mutation.mutateAsync({ status: "approved" });
      toastSuccess("Request approved successfully");
    } catch (error) {
      toastError(
        error instanceof Error ? error.message : "Failed to approve request",
      );
    }
  };

  const handleReject = async (request: UnlockRequest) => {
    const mutation = reviewUnlockRequestMutation(request.id);
    try {
      await mutation.mutateAsync({ status: "rejected" });
      toastSuccess("Request rejected");
    } catch (error) {
      toastError(
        error instanceof Error ? error.message : "Failed to reject request",
      );
    }
  };

  const columns: ColumnDefinition<UnlockRequest>[] = [
    {
      label: "Contest",
      mobile: true,
      render: renderContest,
      width: "1fr",
    },
    {
      label: "Organizer",
      mobile: true,
      render: renderOrganizer,
      width: "1fr",
    },
    {
      label: "Requested",
      mobile: false,
      render: renderTimestamp,
      width: "max-content",
    },
    {
      label: "Status",
      mobile: true,
      render: renderStatus,
      width: "max-content",
    },
    {
      label: "Actions",
      mobile: true,
      render: renderActions,
      width: "max-content",
    },
  ];
</script>

{#snippet renderContest(request: UnlockRequest)}
  {@const contestQuery = getContestQuery(request.contestId)}
  <Link to={`./contests/${request.contestId}`}>
    {contestQuery.data?.name || `Contest ${request.contestId}`}
  </Link>
{/snippet}

{#snippet renderOrganizer(request: UnlockRequest)}
  {@const organizerQuery = getOrganizerQuery(request.organizerId)}
  {organizerQuery.data?.name || `Organizer ${request.organizerId}`}
{/snippet}

{#snippet renderTimestamp(request: UnlockRequest)}
  {new Date(request.createdAt).toLocaleString()}
{/snippet}

{#snippet renderStatus(request: UnlockRequest)}
  <wa-badge variant="warning" size="small">
    {request.status.charAt(0).toUpperCase() + request.status.slice(1)}
  </wa-badge>
{/snippet}

{#snippet renderActions(request: UnlockRequest)}
  {#if request.status === "pending"}
    <div class="actions">
      <wa-button
        size="small"
        variant="success"
        on:click={() => handleApprove(request)}
      >
        Approve
      </wa-button>
      <wa-button size="small" variant="danger" on:click={() => handleReject(request)}>
        Reject
      </wa-button>
    </div>
  {/if}
{/snippet}

<h1>Unlock Requests</h1>

{#if pendingRequestsQuery.isPending}
  <Loader />
{:else if pendingRequestsQuery.isError}
  <p>Error loading unlock requests: {pendingRequestsQuery.error?.message}</p>
{:else if !pendingRequestsQuery.data || pendingRequestsQuery.data.length === 0}
  <div class="empty-state">
    <p>No pending unlock requests. All requests have been reviewed.</p>
  </div>
{:else}
  <Table
    columns={columns}
    data={pendingRequestsQuery.data}
    getId={(request) => request.id}
  />
{/if}

<style>
  h1 {
    margin-block-end: var(--wa-space-l);
  }

  .empty-state {
    padding: var(--wa-space-xl);
    text-align: center;
    color: var(--wa-color-neutral-500);
  }

  .actions {
    display: flex;
    gap: var(--wa-space-xs);
  }
</style>
