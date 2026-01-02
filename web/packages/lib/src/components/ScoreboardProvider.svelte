<script lang="ts">
  import { onDestroy, onMount, type Snippet } from "svelte";
  import { writable, type Writable } from "svelte/store";
  import * as z from "zod/v4";
  import { ApiClient } from "../Api";
  import {
    contenderPublicInfoUpdatedEventSchema,
    contenderScoreUpdatedEventSchema,
    type Score,
    type ScoreboardEntry,
  } from "../models";
  import { getApiUrl } from "../utils";

  interface Props {
    contestId: number;
    hideDisqualified?: boolean;
    children?: Snippet<
      [
        {
          scoreboard: Writable<Map<number, ScoreboardEntry[]>>;
          loading: boolean;
          online: boolean;
        },
      ]
    >;
  }

  let { contestId, hideDisqualified = false, children }: Props = $props();

  let eventSource: EventSource | undefined;
  let loading = $state(true);
  let online = $state(true);

  const contenders: Map<number, ScoreboardEntry> = new Map();
  const pendingUpdates: ((contenders: Map<number, ScoreboardEntry>) => void)[] =
    [];
  const scoreboard = writable<Map<number, ScoreboardEntry[]>>(new Map());

  onMount(async () => {
    setup();
  });

  onDestroy(() => {
    tearDown();
  });

  const handleVisibilityChange = () => {
    switch (document.visibilityState) {
      case "hidden":
        tearDown();
        break;
      case "visible":
        setup();
        break;
    }
  };

  const handleBeforeUnload = () => {
    tearDown();
  };

  const setup = () => {
    eventSource = new EventSource(
      `${getApiUrl()}/contests/${contestId}/events`,
    );

    setupEventHandlers(eventSource);

    eventSource.onerror = () => {
      online = false;
      reset();

      if (eventSource?.readyState === EventSource.CLOSED) {
        setTimeout(() => {
          setup();
        }, 5000);

        return;
      }
    };

    eventSource.onopen = () => {
      online = true;
      initializeStore();
    };
  };

  const reset = () => {
    loading = true;
    contenders.clear();
    $scoreboard = new Map();
  };

  const tearDown = () => {
    eventSource?.close();
    eventSource = undefined;

    reset();
  };

  const initializeStore = async () => {
    const entries = await ApiClient.getInstance().getScoreboard(contestId);

    for (const entry of entries) {
      contenders.set(entry.contenderId, entry);
    }

    while (pendingUpdates.length > 0) {
      const handler = pendingUpdates.shift();
      handler?.(contenders);
    }

    rebuildStore();
    loading = false;
  };

  const rebuildStore = () => {
    const results = new Map<number, ScoreboardEntry[]>();
    const withoutScore: ScoreboardEntry[] = [];

    for (const contender of contenders.values()) {
      let classEntries = results.get(contender.compClassId);

      if (classEntries === undefined) {
        classEntries = [];
        results.set(contender.compClassId, classEntries);
      }

      if (hideDisqualified && contender.disqualified) {
        continue;
      }

      if (!contender.score) {
        withoutScore.push(contender);
      } else {
        classEntries.push(contender);
      }
    }

    if (withoutScore.length) {
      for (const [compClassId, classEntries] of results.entries()) {
        const maxRankOrder = Math.max(
          ...classEntries.map(({ score }) => score?.rankOrder ?? -1),
          -1,
        );

        let rankOrder = maxRankOrder + 1;

        for (const contender of withoutScore) {
          if (contender.compClassId !== compClassId) {
            continue;
          }

          contender.score = {
            contenderId: contender.contenderId,
            score: 0,
            placement: 0,
            finalist: false,
            rankOrder,
            timestamp: new Date(0),
          };

          classEntries.push(contender);

          rankOrder++;
        }
      }
    }

    $scoreboard = results;
  };

  const queueEventHandler = (
    handler: (contenders: Map<number, ScoreboardEntry>) => void,
  ) => {
    if (loading) {
      pendingUpdates.push(handler);
    } else {
      handler(contenders);
      rebuildStore();
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
        }

        contender.compClassId = event.compClassId;
        contender.name = event.name;
        contender.withdrawnFromFinals = event.withdrawnFromFinals;
        contender.disqualified = event.disqualified;

        contenders.set(event.contenderId, { ...contender });
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
          }

          const score: Score = {
            contenderId: event.contenderId,
            score: event.score,
            placement: event.placement,
            finalist: event.finalist,
            rankOrder: event.rankOrder,
            timestamp: event.timestamp,
          };

          contender.score = score;
          contenders.set(event.contenderId, { ...contender });
        });
      }
    });
  };

  const createEmptyEntry = (contenderId: number): ScoreboardEntry => ({
    contenderId: contenderId,
    compClassId: 0,
    name: "",
    withdrawnFromFinals: false,
    disqualified: false,
  });
</script>

<svelte:window
  onbeforeunload={handleBeforeUnload}
  onvisibilitychange={handleVisibilityChange}
/>

{@render children?.({ scoreboard, loading, online })}
