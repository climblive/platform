<script lang="ts" module>
  export const doSubmit = (
    mutation: ReturnType<typeof patchCompetitionMutation>,
    patch: CompetitionPatch,
  ) => {
    mutation.mutate(patch, {
      onError: () => toastUnexpectedError("Failed to update rules."),
    });
  };
</script>

<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/radio/radio.js";
  import "@awesome.me/webawesome/dist/components/radio-group/radio-group.js";
  import type WaRadioGroup from "@awesome.me/webawesome/dist/components/radio-group/radio-group.js";
  import { SaveIndicator } from "@climblive/lib/components";
  import type { Competition, CompetitionPatch } from "@climblive/lib/models";
  import { patchCompetitionMutation } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import RuleOptionCard from "./RuleOptionCard.svelte";
  import Finalists from "./rules/Finalists.svelte";
  import PooledPoints from "./rules/PooledPoints.svelte";
  import MaxAttempts from "./rules/MaxAttempts.svelte";
  import PointDeduction from "./rules/PointDeduction.svelte";
  import ProblemLimit from "./rules/ProblemLimit.svelte";

  interface Props {
    competition: Competition;
  }

  const { competition }: Props = $props();

  const patchCompetition = $derived(patchCompetitionMutation(competition.id));

  const handleUsePointsChange = (usePoints: boolean) => {
    const patch: CompetitionPatch = {
      usePoints,
    };

    if (!usePoints) {
      patch.qualifyingProblems = 0;
      patch.pooledPoints = false;
      patch.maxAttempts = 0;
      patch.pointDeduction = 0;
    }

    patchCompetition.mutate(patch, {
      onError: () => toastUnexpectedError("Failed to update rules."),
    });
  };
</script>

<h3>Ranking method</h3>
<section>
  <RuleOptionCard
    title="Points"
    description="Competitors are ranked based on the total points scored across all problems."
    tags={["Recommended"]}
  >
    {#snippet header()}
      <wa-radio
        onclick={() => handleUsePointsChange(true)}
        size="s"
        checked={competition.usePoints ? true : undefined}
      ></wa-radio>
    {/snippet}
  </RuleOptionCard>

  <RuleOptionCard
    title="Attempts"
    description="Competitors are ranked by tops and zones, with attempts used to break ties."
  >
    {#snippet header()}
      <wa-radio
        onclick={() => handleUsePointsChange(false)}
        size="s"
        checked={!competition.usePoints ? true : undefined}
      ></wa-radio>
    {/snippet}
  </RuleOptionCard>
</section>

{#if !competition.usePoints}
  <div class="attempt-priority">
    <wa-radio-group
      label="Tie-breaking priority"
      hint="Tops always rank first. This setting determines how ties in tops are resolved."
      size="s"
      value={competition.favorTopAttempts ? "attempts" : "zones"}
      disabled={patchCompetition.isPending}
      onchange={(event: Event) =>
        doSubmit(patchCompetition, {
          favorTopAttempts:
            (event.currentTarget as WaRadioGroup).value === "attempts",
        })}
    >
      <wa-radio value="zones">
        Favor zones over attempts (IFSC Boulder 2018–2024)
      </wa-radio>
      <wa-radio value="attempts">
        Favor attempts to tops (IFSC Boulder through 2017)
      </wa-radio>
    </wa-radio-group>
    <p>
      These historical formats used one zone. With two zones, zones-first
      compares both zone counts before attempts; the older format compares
      attempts at each level before moving to the next zone.
    </p>
    {#if patchCompetition.isSuccess}
      <SaveIndicator />
    {/if}
  </div>
{/if}

<h3>Options</h3>
<section>
  <Finalists {competition} />

  <PooledPoints {competition} />

  <ProblemLimit {competition} />

  <MaxAttempts {competition} />

  <PointDeduction {competition} />
</section>

<style>
  .attempt-priority {
    margin-block-start: var(--wa-space-m);

    p {
      font-size: var(--wa-font-size-s);
      color: var(--wa-color-text-quiet);
    }
  }

  section {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--wa-space-m);
  }

  @media screen and (max-width: 768px) {
    section {
      grid-template-columns: 1fr;
    }
  }
</style>
