<script lang="ts">
  import { ApiClient } from "@climblive/lib";
  import configData from "@climblive/lib/config.json";
  import type {
    ContenderPublicInfoUpdatedEvent,
    ContenderScoreUpdatedEvent,
    ScoreboardEntry,
  } from "@climblive/lib/models";
  import { onDestroy, onMount, setContext } from "svelte";
  import { writable } from "svelte/store";

  export let contestId: number;

  let eventSource: EventSource | undefined;
  let initialized = false;

  const contenders: Map<number, ScoreboardEntry> = new Map();
  const pendingUpdates: ((contenders: Map<number, ScoreboardEntry>) => void)[] =
    [];

  const resultsStore = writable<Map<number, ScoreboardEntry[]>>(new Map());

  setContext("scoreboard", resultsStore);

  onMount(async () => {
    const entries = await ApiClient.getInstance().getScoreboard(contestId);

    for (const entry of entries) {
      contenders.set(entry.contenderId, entry);
    }

    pendingUpdates.forEach((handler) => handler(contenders));

    rebuildStore();
    initialized = true;
  });

  const rebuildStore = () => {
    const results = new Map<number, ScoreboardEntry[]>();

    for (const contender of contenders.values()) {
      let classEntries: ScoreboardEntry[];

      if (results.has(contender.compClassId)) {
        classEntries = results.get(contender.compClassId)!;
      } else {
        classEntries = [];
        results.set(contender.compClassId, classEntries);
      }

      classEntries.push(contender);
    }

    $resultsStore = results;
  };

  const queueEventHandler = (
    handler: (contenders: Map<number, ScoreboardEntry>) => void,
  ) => {
    if (initialized) {
      handler(contenders);
      rebuildStore();
    } else {
      pendingUpdates.push(handler);
    }
  };

  onMount(() => {
    eventSource = new EventSource(
      `${configData.API_URL}/contests/${contestId}/events`,
    );

    eventSource.addEventListener("CONTENDER_PUBLIC_INFO_UPDATED", (e) => {
      const event = JSON.parse(e.data) as ContenderPublicInfoUpdatedEvent;

      queueEventHandler((contenders: Map<number, ScoreboardEntry>) => {
        const contender = contenders.get(event.contenderId);
        if (!contender) {
          return;
        }

        contender.compClassId = event.compClassId;
        contender.publicName = event.publicName;
        contender.clubName = event.clubName;
        contender.withdrawnFromFinals = event.withdrawnFromFinals;
        contender.disqualified = event.disqualified;
      });
    });

    eventSource.addEventListener("CONTENDER_SCORE_UPDATED", (e) => {
      const event = JSON.parse(e.data) as ContenderScoreUpdatedEvent;

      queueEventHandler((contenders: Map<number, ScoreboardEntry>) => {
        const contender = contenders.get(event.contenderId);
        if (!contender) {
          return;
        }

        contender.score = event.score;
        contender.placement = event.placement;
        contender.finalist = event.finalist;
      });
    });
  });

  onDestroy(() => {
    eventSource?.close();
    eventSource = undefined;
  });
</script>

<slot />
