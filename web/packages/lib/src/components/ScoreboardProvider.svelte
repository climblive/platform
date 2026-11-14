<script lang="ts">
  import { z } from "@climblive/lib/utils";
  import { onDestroy, onMount, type Snippet } from "svelte";
  import { writable, type Writable } from "svelte/store";
  import { ApiClient } from "../Api";
  import {
    contenderPublicInfoUpdatedEventSchema,
    contenderScoreUpdatedEventSchema,
    type EventFilter,
    type Score,
    type ScoreboardEntry,
  } from "../models";
  import { NO_SCORE } from "../models/score";
  import { useEvents } from "./EventsProvider.svelte";

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

  const { subscribeEvents, unsubscribeEvents } = useEvents();
  let subscriptionId: symbol | undefined;
  let loading = $state(true);
  let online = $state(true);
  let generation = 0;

  const contenders: Map<number, ScoreboardEntry> = new Map();
  const pendingUpdates: ((contenders: Map<number, ScoreboardEntry>) => void)[] =
    [];
  const scoreboard = writable<Map<number, ScoreboardEntry[]>>(new Map());

  onMount(() => {
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
    if (subscriptionId !== undefined) {
      return;
    }

    const filters: EventFilter[] = [
      {
        contestId,
        eventTypes: {
          CONTENDER_PUBLIC_INFO_UPDATED: {},
          "[]CONTENDER_SCORE_UPDATED": {},
        },
      },
    ];

    subscriptionId = subscribeEvents(filters, handleEvent, {
      onError: () => {
        online = false;
        reset();
      },
      onOpen: () => {
        online = true;
        reset();
        initializeStore();
      },
    });
  };

  const reset = () => {
    generation++;
    loading = true;
    contenders.clear();
    pendingUpdates.length = 0;
    $scoreboard = new Map();
  };

  const tearDown = () => {
    if (subscriptionId !== undefined) {
      unsubscribeEvents(subscriptionId);
      subscriptionId = undefined;
    }

    reset();
  };

  const initializeStore = async () => {
    const currentGeneration = generation;
    const entries = await ApiClient.getInstance().getScoreboard(contestId);

    if (currentGeneration !== generation) {
      return;
    }

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
            score: NO_SCORE,
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

  const handleEvent = (eventName: string, e: MessageEvent) => {
    if (eventName === "CONTENDER_PUBLIC_INFO_UPDATED") {
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
        contender.scrubbedAt = event.scrubbedAt;

        contenders.set(event.contenderId, { ...contender });
      });
    }

    if (eventName === "[]CONTENDER_SCORE_UPDATED") {
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
    }
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
