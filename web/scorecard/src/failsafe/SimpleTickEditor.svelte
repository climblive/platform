<script lang="ts">
  import type { Problem, Tick } from "@climblive/lib/models";
  import { deleteTickMutation, putTickMutation } from "@climblive/lib/queries";

  type Props = {
    problem: Problem;
    tick?: Tick;
    contenderId: number;
  };

  const { problem, tick, contenderId }: Props = $props();

  const putTick = $derived(putTickMutation(contenderId));
  let latestLocalRevision = $state(0);
  const deleteTick = $derived(deleteTickMutation());

  const addTick = (type: "zone1" | "zone2" | "top" | "flash") => () => {
    const attempts = type === "flash" ? 1 : 999;
    latestLocalRevision =
      Math.max(latestLocalRevision, tick?.revision ?? 0) + 1;

    const nextTick: Omit<Tick, "id" | "timestamp"> = {
      revision: latestLocalRevision,
      problemId: problem.id,
      top: false,
      zone2: false,
      zone1: false,
      attemptsTop: attempts,
      attemptsZone2: attempts,
      attemptsZone1: attempts,
    };

    switch (type) {
      case "flash":
      case "top":
        nextTick.top = true;
        nextTick.zone2 = true;
        nextTick.zone1 = true;
        break;
      case "zone2":
        nextTick.zone2 = true;
        nextTick.zone1 = true;
        break;
      case "zone1":
        nextTick.zone1 = true;
        break;
    }

    putTick.mutate(nextTick, {
      onError: () => window.alert("Failed to update ascent."),
    });
  };

  const removeTick = () => {
    if (tick?.id) {
      deleteTick.mutate(tick.id, {
        onError: () => window.alert("Failed to remove ascent."),
      });
    }
  };
</script>

<div>
  {#if tick}
    <button type="button" onclick={removeTick} disabled={deleteTick.isPending}
      >Unsend</button
    >
  {:else}
    {#if problem.zone1Enabled}
      <button
        type="button"
        onclick={addTick("zone1")}
        disabled={putTick.isPending}>Zone 1</button
      >
    {/if}
    {#if problem.zone2Enabled}
      <button
        type="button"
        onclick={addTick("zone2")}
        disabled={putTick.isPending}>Zone 2</button
      >
    {/if}
    <button type="button" onclick={addTick("top")} disabled={putTick.isPending}
      >Top</button
    >
    <button
      type="button"
      onclick={addTick("flash")}
      disabled={putTick.isPending}>Flash</button
    >
  {/if}
</div>
