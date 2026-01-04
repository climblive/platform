<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import { toastError, toastSuccess } from "@climblive/lib/utils";
  import type { UnlockRequest } from "@climblive/lib/models";
  import {
    getContestQuery,
    getOrganizerQuery,
    getPendingUnlockRequestsQuery,
    reviewUnlockRequestMutation,
  } from "@climblive/lib/queries";

  const pendingRequestsQuery = getPendingUnlockRequestsQuery();

  // Web component types for wa-dialog
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let approveDialog: any = $state();
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let rejectDialog: any = $state();
  let selectedRequest: UnlockRequest | null = $state(null);

  const approveMutation = $derived(
    selectedRequest
      ? reviewUnlockRequestMutation((selectedRequest as UnlockRequest).id)
      : undefined,
  ) as ReturnType<typeof reviewUnlockRequestMutation> | undefined;

  const rejectMutation = $derived(
    selectedRequest
      ? reviewUnlockRequestMutation((selectedRequest as UnlockRequest).id)
      : undefined,
  ) as ReturnType<typeof reviewUnlockRequestMutation> | undefined;

  const handleApproveClick = (request: UnlockRequest) => {
    selectedRequest = request;
    approveDialog?.show();
  };

  const handleRejectClick = (request: UnlockRequest) => {
    selectedRequest = request;
    rejectDialog?.show();
  };

  const handleApprove = async () => {
    if (!approveMutation || !selectedRequest) return;

    try {
      await approveMutation.mutateAsync({ status: "approved" });
      toastSuccess("Request approved successfully");
      approveDialog?.hide();
      selectedRequest = null;
    } catch (error) {
      toastError(
        error instanceof Error ? error.message : "Failed to approve request",
      );
    }
  };

  const handleReject = async () => {
    if (!rejectMutation || !selectedRequest) return;

    try {
      await rejectMutation.mutateAsync({ status: "rejected" });
      toastSuccess("Request rejected");
      rejectDialog?.hide();
      selectedRequest = null;
    } catch (error) {
      toastError(
        error instanceof Error ? error.message : "Failed to reject request",
      );
    }
  };

  const selectedContestQuery = $derived(
    selectedRequest
      ? getContestQuery((selectedRequest as UnlockRequest).contestId)
      : undefined,
  );
  const selectedContestName = $derived(
    selectedContestQuery?.data?.name ||
      `Contest ${(selectedRequest as UnlockRequest | null)?.contestId || ""}`,
  );
</script>

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
  <table>
    <thead>
      <tr>
        <th>Contest</th>
        <th>Organizer</th>
        <th>Requested</th>
        <th>Status</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      {#each pendingRequestsQuery.data as request (request.id)}
        {@const contestQuery = getContestQuery(request.contestId)}
        {@const organizerQuery = getOrganizerQuery(request.organizerId)}
        <tr>
          <td>
            <a href="/contests/{request.contestId}">
              {contestQuery.data?.name || `Contest ${request.contestId}`}
            </a>
          </td>
          <td>
            {organizerQuery.data?.name || `Organizer ${request.organizerId}`}
          </td>
          <td>
            {new Date(request.createdAt).toLocaleString()}
          </td>
          <td>
            <wa-badge variant="warning" size="small">
              {request.status.charAt(0).toUpperCase() + request.status.slice(1)}
            </wa-badge>
          </td>
          <td>
            {#if request.status === "pending"}
              <div class="actions">
                <wa-button
                  size="small"
                  variant="success"
                  onclick={() => handleApproveClick(request)}
                >
                  Approve
                </wa-button>
                <wa-button
                  size="small"
                  variant="danger"
                  onclick={() => handleRejectClick(request)}
                >
                  Reject
                </wa-button>
              </div>
            {/if}
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
{/if}

<wa-dialog bind:this={approveDialog} label="Approve unlock request">
  <p>
    Approve unlock request for <strong>{selectedContestName}</strong>?
  </p>
  <p>This will unlock the contest to its full capacity of 500 contenders.</p>
  <wa-button slot="footer" variant="text" onclick={() => approveDialog?.hide()}>
    Cancel
  </wa-button>
  <wa-button
    slot="footer"
    variant="primary"
    onclick={handleApprove}
    loading={approveMutation?.isPending}
    disabled={approveMutation?.isPending}
  >
    Approve
  </wa-button>
</wa-dialog>

<wa-dialog bind:this={rejectDialog} label="Reject unlock request">
  <p>
    Reject unlock request for <strong>{selectedContestName}</strong>?
  </p>
  <p>The contest will remain in evaluation mode with a 10 contender limit.</p>
  <wa-button slot="footer" variant="text" onclick={() => rejectDialog?.hide()}>
    Cancel
  </wa-button>
  <wa-button
    slot="footer"
    variant="danger"
    onclick={handleReject}
    loading={rejectMutation?.isPending}
    disabled={rejectMutation?.isPending}
  >
    Reject
  </wa-button>
</wa-dialog>

<style>
  h1 {
    margin-block-end: var(--wa-space-l);
  }

  .empty-state {
    padding: var(--wa-space-xl);
    text-align: center;
    color: var(--wa-color-neutral-500);
  }

  table {
    width: 100%;
    border-collapse: collapse;
    margin-block-start: var(--wa-space-m);
  }

  th,
  td {
    padding: var(--wa-space-m);
    text-align: left;
    border-bottom: 1px solid var(--wa-color-neutral-200);
  }

  th {
    font-weight: 600;
    color: var(--wa-color-neutral-700);
  }

  .actions {
    display: flex;
    gap: var(--wa-space-xs);
  }

  a {
    color: var(--wa-color-primary-600);
    text-decoration: none;
  }

  a:hover {
    text-decoration: underline;
  }
</style>
