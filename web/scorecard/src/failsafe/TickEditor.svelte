<script lang="ts">
  import type { Problem, Tick } from "@climblive/lib/models";
  import { deleteTickMutation, putTickMutation } from "@climblive/lib/queries";
  import {
    buildTick,
    TickMutator,
    type Feature,
  } from "../utils/tickMutator.svelte";

  type Props = {
    problem: Problem;
    tick?: Tick;
    contenderId: number;
  };

  const { problem, tick, contenderId }: Props = $props();

  const putTick = $derived(putTickMutation(contenderId));
  const deleteTick = $derived(deleteTickMutation());
  let tickMutator = $derived(TickMutator.from(problem, tick));
  let latestLocalRevision = $state(0);

  const saveTick = async () => {
    const revision = Math.max(latestLocalRevision, tick?.revision ?? 0) + 1;
    latestLocalRevision = revision;

    try {
      await putTick.mutateAsync({
        ...buildTick(problem.id, tickMutator),
        revision,
      });
    } catch {
      tickMutator = TickMutator.from(problem, tick);
      window.alert("Failed to update ascent. Changes reverted.");
    }
  };

  const handleTick = (checked: boolean, feature: Feature) => {
    if (checked) {
      tickMutator.reachFeature(feature);
    } else {
      tickMutator.unreachFeature(feature);
    }

    saveTick();
  };

  const removeTick = () => {
    if (tick?.id) {
      deleteTick.mutate(tick.id, {
        onError: () => window.alert("Failed to remove ascent."),
      });
    }
  };

  const isPending = $derived(putTick.isPending || deleteTick.isPending);
</script>

<div>
  <button
    type="button"
    aria-label="Subtract failed attempt"
    disabled={!tickMutator.canSubtractAttempt() || isPending}
    onclick={() => {
      tickMutator.subtractAttempt();
      saveTick();
    }}>−</button
  >
  <span role="status" aria-label="Attempts">
    {tickMutator.attempts}
    {tickMutator.attempts === 1 ? "attempt" : "attempts"}
  </span>
  <button
    type="button"
    disabled={!tickMutator.canAddAttempt() || isPending}
    onclick={() => {
      tickMutator.addAttempt();
      saveTick();
    }}>Add failed attempt</button
  >
</div>

{#each [...tickMutator.features].reverse() as feature (feature)}
  {@const reachedAt = tickMutator.reachedFeatures.get(feature)}
  {@const attempts = reachedAt ?? tickMutator.attempts + 1}
  <div>
    <input
      type="checkbox"
      checked={reachedAt !== undefined}
      disabled={isPending ||
        (reachedAt === undefined && !tickMutator.canAddAttempt())}
      onchange={(event) => handleTick(event.currentTarget.checked, feature)}
    />
    {feature === "top" ? "Top" : feature === "zone2" ? "Zone 2" : "Zone 1"}
    {#if attempts <= 999}
      in {attempts} {attempts === 1 ? "attempt" : "attempts"}
    {/if}
  </div>
{/each}

{#if tick}
  <button
    type="button"
    onclick={removeTick}
    disabled={putTick.isPending || deleteTick.isPending}>Remove</button
  >
{/if}
