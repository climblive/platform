<script lang="ts">
  import { ApiClient } from "@climblive/lib";
  import configData from "@climblive/lib/config.json";
  import {
    contenderPublicInfoUpdatedEventSchema,
    contenderScoreUpdatedEventSchema,
    type ScoreboardEntry,
  } from "@climblive/lib/models";
  import { onDestroy, onMount, setContext } from "svelte";
  import { writable } from "svelte/store";
  import * as z from "zod";

  export let contestId: number;

  let eventSource: EventSource | undefined;
  let initialized = false;

  const contenders: Map<number, ScoreboardEntry> = new Map();
  const pendingUpdates: ((contenders: Map<number, ScoreboardEntry>) => void)[] =
    [];
  const scoreboardStore = writable<Map<number, ScoreboardEntry[]>>(new Map());

  setContext("scoreboard", scoreboardStore);

  onMount(async () => {
    eventSource = new EventSource(
      `${configData.API_URL}/contests/${contestId}/events`,
    );

    setupEventHandlers(eventSource);

    eventSource.onerror = (e) => {
      // eslint-disable-next-line no-console
      console.error(e);

      initialized = false;
      contenders.clear();
      $scoreboardStore = new Map();
    };

    eventSource.onopen = () => {
      initializeStore();
    };
  });

  onDestroy(() => {
    eventSource?.close();
    eventSource = undefined;
  });

  const initializeStore = async () => {
    const entries = await ApiClient.getInstance().getScoreboard(contestId);

    for (const entry of entries) {
      contenders.set(entry.contenderId, entry);
    }

    pendingUpdates.forEach((handler) => handler(contenders));

    rebuildStore();
    initialized = true;
  };

  const rebuildStore = () => {
    const results = new Map<number, ScoreboardEntry[]>();

    for (const contender of contenders.values()) {
      if (!contender.scoreUpdated) {
        continue;
      }

      let classEntries = results.get(contender.compClassId);

      if (classEntries === undefined) {
        classEntries = [];
        results.set(contender.compClassId, classEntries);
      }

      classEntries.push(contender);
    }

    $scoreboardStore = results;
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

  const setupEventHandlers = (eventSource: EventSource) => {
    eventSource.addEventListener("CONTENDER_PUBLIC_INFO_UPDATED", (e) => {
      const event = contenderPublicInfoUpdatedEventSchema.parse(
        JSON.parse(e.data),
      );

      queueEventHandler((contenders: Map<number, ScoreboardEntry>) => {
        let contender = contenders.get(event.contenderId);
        if (!contender) {
          contender = createEmptyEntry(event.contenderId);
          contenders.set(event.contenderId, contender);
        }

        contender.compClassId = event.compClassId;
        contender.publicName = event.publicName;
        contender.clubName = event.clubName;
        contender.withdrawnFromFinals = event.withdrawnFromFinals;
        contender.disqualified = event.disqualified;
      });
    });

    eventSource.addEventListener("[]CONTENDER_SCORE_UPDATED", (e) => {
      const events = z
        .array(contenderScoreUpdatedEventSchema)
        .parse(JSON.parse(e.data));

      for (const event of events) {
        queueEventHandler((contenders: Map<number, ScoreboardEntry>) => {
          let contender = contenders.get(event.contenderId);
          if (!contender) {
            contender = createEmptyEntry(event.contenderId);
            contenders.set(event.contenderId, contender);
          }

          contender.score = event.score;
          contender.placement = event.placement;
          contender.finalist = event.finalist;
          contender.rankOrder = event.rankOrder;
          contender.scoreUpdated = event.timestamp;
        });
      }
    });
  };

  const createEmptyEntry = (contenderId: number): ScoreboardEntry => ({
    contenderId: contenderId,
    compClassId: 0,
    withdrawnFromFinals: false,
    disqualified: false,
    score: 0,
    rankOrder: 0,
    finalist: false,
  });
</script>

<slot />
