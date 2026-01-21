<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/progress-bar/progress-bar.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import { EmptyState } from "@climblive/lib/components";
  import { type ProblemTemplate } from "@climblive/lib/models";
  import {
    createProblemMutation,
    getAllContestsQuery,
    getProblemsQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";

  interface Props {
    contestId: number;
    open: boolean;
  }

  let { contestId, open = $bindable() }: Props = $props();

  let dialog: WaDialog | undefined = $state();
  let selectedContestId: number | undefined = $state();
  let isCopying = $state(false);
  let copyProgress = $state(0);
  let copyCompleted = $state(false);

  const contestsQuery = $derived(getAllContestsQuery());

  const availableContests = $derived.by(() => {
    if (!contestsQuery.data) {
      return [];
    }

    return contestsQuery.data
      .filter(({ id, archived }) => id !== contestId && !archived)
      .sort((a, b) => b.created.getTime() - a.created.getTime());
  });

  const selectedContest = $derived(
    availableContests.find((c) => c.id === selectedContestId),
  );

  const selectedProblemsQuery = $derived(
    selectedContest ? getProblemsQuery(selectedContest.id) : undefined,
  );

  const problems = $derived(selectedProblemsQuery?.data);

  const createMutation = $derived(createProblemMutation(contestId));

  const problemsSummary = $derived.by(() => {
    if (!problems || problems.length === 0) {
      return null;
    }

    const points = problems.map((p) => p.pointsTop);
    const minPoints = Math.min(...points);
    const maxPoints = Math.max(...points);

    return {
      count: problems.length,
      minPoints,
      maxPoints,
    };
  });

  const handleClose = () => {
    open = false;
  };

  const handleContestChange = (e: Event) => {
    const select = e.target as WaSelect;
    selectedContestId = select.value ? Number(select.value) : undefined;
  };

  const handleCopy = async () => {
    if (!problems || problems.length === 0) {
      return;
    }

    isCopying = true;
    copyProgress = 0;
    copyCompleted = false;
    let failCount = 0;

    for (let i = 0; i < problems.length; i++) {
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const { id, contestId, ...rest } = problems[i];

      const template: ProblemTemplate = {
        ...rest,
      };

      try {
        await createMutation.mutateAsync(template);
      } catch {
        failCount++;
      }

      copyProgress = ((i + 1) / problems.length) * 100;
    }

    isCopying = false;
    copyCompleted = true;

    if (failCount > 0) {
      toastError(
        `Failed to copy ${failCount} problem${failCount > 1 ? "s" : ""}`,
      );
    }
  };
</script>

<wa-dialog bind:this={dialog} label="Copy problems" {open}>
  {#if contestsQuery.isPending}
    <Loader />
  {:else if availableContests.length === 0}
    <EmptyState
      title="No contests available"
      description="There are no other contests to copy problems from."
    />
  {:else}
    {#if !isCopying && !copyCompleted}
      <wa-select
        label="Select contest"
        placeholder="Select a contest to copy from"
        onchange={handleContestChange}
      >
        {#each availableContests as contest (contest.id)}
          <wa-option value={contest.id}>{contest.name}</wa-option>
        {/each}
      </wa-select>
    {/if}

    {#if selectedContest}
      {#if selectedProblemsQuery?.isPending}
        <Loader />
      {:else if problemsSummary}
        {#if isCopying || copyCompleted}
          <wa-progress-bar value={copyProgress}></wa-progress-bar>
        {:else}
          <p class="summary">
            This contest has {problemsSummary.count} problem{problemsSummary.count !==
            1
              ? "s"
              : ""} with point values ranging from {problemsSummary.minPoints} to
            {problemsSummary.maxPoints}.
          </p>
        {/if}
      {:else}
        <EmptyState
          title="No problems found"
          description="This contest has no problems to copy."
        />
      {/if}
    {/if}
  {/if}

  {#if copyCompleted}
    <wa-button slot="footer" variant="success" onclick={handleClose}>
      <wa-icon name="check" slot="start"></wa-icon>
      Close
    </wa-button>
  {:else}
    <wa-button slot="footer" appearance="plain" onclick={handleClose}>
      Cancel
    </wa-button>
    {#if selectedContest}
      <wa-button
        slot="footer"
        variant="accent"
        onclick={handleCopy}
        disabled={!problems || problems.length === 0 || isCopying}
        loading={isCopying}
      >
        Copy {problems?.length || 0} problem{problems?.length !== 1 ? "s" : ""}
        <wa-icon slot="start" name="copy"></wa-icon>
      </wa-button>
    {/if}
  {/if}
</wa-dialog>

<style>
  wa-dialog::part(body) {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }

  .summary {
    margin: 0;
  }
</style>
