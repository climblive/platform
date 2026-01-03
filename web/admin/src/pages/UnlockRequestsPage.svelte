<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import RelativeTime from "@/components/RelativeTime.svelte";
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import {
    EmptyState,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import type { UnlockRequest } from "@climblive/lib/models";
  import {
    getContestQuery,
    getOrganizerQuery,
    getPendingUnlockRequestsQuery,
    reviewUnlockRequestMutation,
  } from "@climblive/lib/queries";
  import { toastError, toastSuccess } from "@climblive/lib/utils/errors";
  import { Link } from "svelte-routing";

  const pendingRequestsQuery = getPendingUnlockRequestsQuery();

  let approveDialog: any = $state();
  let rejectDialog: any = $state();
  let selectedRequest: UnlockRequest | null = $state(null);

  const approveMutation = $derived(
    selectedRequest
      ? reviewUnlockRequestMutation(selectedRequest.id)
      : undefined,
  );

  const rejectMutation = $derived(
    selectedRequest
      ? reviewUnlockRequestMutation(selectedRequest.id)
      : undefined,
  );

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
      toastError(error);
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
      toastError(error);
    }
  };

  // Helper component for contest name
  function ContestName(props: { contestId: number }) {
    const contestQuery = getContestQuery(props.contestId);
    return contestQuery.data?.name || `Contest ${props.contestId}`;
  }

  // Helper component for organizer name
  function OrganizerName(props: { organizerId: number }) {
    const organizerQuery = getOrganizerQuery(props.organizerId);
    return organizerQuery.data?.name || `Organizer ${props.organizerId}`;
  }

  const columns: ColumnDefinition<UnlockRequest>[] = [
    {
      label: "Contest",
      prop: "contestId",
      cellRenderer: (request) => (
        <Link to={`/contests/${request.contestId}`}>
          <ContestName contestId={request.contestId} />
        </Link>
      ),
    },
    {
      label: "Organizer",
      prop: "organizerId",
      cellRenderer: (request) => (
        <span>
          <OrganizerName organizerId={request.organizerId} />
        </span>
      ),
    },
    {
      label: "Requested",
      prop: "createdAt",
      cellRenderer: (request) => <RelativeTime time={request.createdAt} />,
    },
    {
      label: "Status",
      prop: "status",
      cellRenderer: (request) => {
        const statusMap = {
          pending: { variant: "warning", text: "Pending" },
          approved: { variant: "success", text: "Approved" },
          rejected: { variant: "danger", text: "Rejected" },
        };
        const status = statusMap[request.status] || statusMap.pending;
        return (
          <wa-badge variant={status.variant} size="small">
            {status.text}
          </wa-badge>
        );
      },
    },
    {
      label: "Actions",
      prop: "id",
      cellRenderer: (request) =>
        request.status === "pending" ? (
          <div style="display: flex; gap: var(--wa-space-xs);">
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
        ) : null,
    },
  ];

  // Fetch contest name for dialog
  const selectedContestQuery = $derived(
    selectedRequest ? getContestQuery(selectedRequest.contestId) : undefined,
  );
  const selectedContestName = $derived(
    selectedContestQuery?.data?.name ||
      `Contest ${selectedRequest?.contestId || ""}`,
  );
</script>

<h1>Unlock Requests</h1>

{#if pendingRequestsQuery.isPending}
  <Loader />
{:else if pendingRequestsQuery.isError}
  <p>Error loading unlock requests: {pendingRequestsQuery.error?.message}</p>
{:else if !pendingRequestsQuery.data || pendingRequestsQuery.data.length === 0}
  <EmptyState
    title="No pending unlock requests"
    message="All unlock requests have been reviewed."
  />
{:else}
  <Table data={pendingRequestsQuery.data} {columns} />
{/if}

<wa-dialog bind:this={approveDialog} label="Approve unlock request">
  <p>
    Approve unlock request for <strong>{selectedContestName}</strong>?
  </p>
  <p>
    This will unlock the contest to its full capacity of 500 contenders.
  </p>
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
</style>
