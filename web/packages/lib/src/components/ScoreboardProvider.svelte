<script lang="ts">
  import { z } from "@climblive/lib/utils";
  import { onDestroy, onMount, type Snippet } from "svelte";
  import { writable, type Writable } from "svelte/store";
  import { ApiClient } from "../Api";
  import {
    competitorPublicInfoUpdatedEventSchema,
    competitorResultUpdatedEventSchema,
    type EventFilter,
    type Result,
    type ScoreboardEntry,
  } from "../models";
  import { NO_SCORE } from "../models/score";
  import { useEvents } from "./EventsProvider.svelte";

  interface Props {
    competitionId: string;
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

  let { competitionId, hideDisqualified = false, children }: Props = $props();

  const { subscribeEvents, unsubscribeEvents } = useEvents();
  let subscriptionId: symbol | undefined;
  let loading = $state(true);
  let online = $state(true);
  let generation = 0;

  const competitors: Map<number, ScoreboardEntry> = new Map();
  const pendingUpdates: ((
    competitors: Map<number, ScoreboardEntry>,
  ) => void)[] = [];
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
        competitionId,
        eventTypes: {
          COMPETITOR_PUBLIC_INFO_UPDATED: {},
          "[]COMPETITOR_RESULT_UPDATED": {},
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
    competitors.clear();
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
    const entries = await ApiClient.getInstance().getScoreboard(competitionId);

    if (currentGeneration !== generation) {
      return;
    }

    for (const entry of entries) {
      competitors.set(entry.competitorId, entry);
    }

    while (pendingUpdates.length > 0) {
      const handler = pendingUpdates.shift();
      handler?.(competitors);
    }

    rebuildStore();
    loading = false;
  };

  const rebuildStore = () => {
    const results = new Map<number, ScoreboardEntry[]>();
    const withoutResult: ScoreboardEntry[] = [];

    for (const competitor of competitors.values()) {
      let categoryEntries = results.get(competitor.categoryId);

      if (categoryEntries === undefined) {
        categoryEntries = [];
        results.set(competitor.categoryId, categoryEntries);
      }

      if (hideDisqualified && competitor.disqualified) {
        continue;
      }

      if (!competitor.result) {
        withoutResult.push(competitor);
      } else {
        categoryEntries.push(competitor);
      }
    }

    if (withoutResult.length) {
      for (const [categoryId, categoryEntries] of results.entries()) {
        const maxRankOrder = Math.max(
          ...categoryEntries.map(({ result }) => result?.rankOrder ?? -1),
          -1,
        );

        let rankOrder = maxRankOrder + 1;

        for (const competitor of withoutResult) {
          if (competitor.categoryId !== categoryId) {
            continue;
          }

          competitor.result = {
            competitorId: competitor.competitorId,
            score: NO_SCORE,
            placement: 0,
            finalist: false,
            rankOrder,
            timestamp: new Date(0),
          };

          categoryEntries.push(competitor);

          rankOrder++;
        }
      }
    }

    $scoreboard = results;
  };

  const queueEventHandler = (
    handler: (competitors: Map<number, ScoreboardEntry>) => void,
  ) => {
    if (loading) {
      pendingUpdates.push(handler);
    } else {
      handler(competitors);
      rebuildStore();
    }
  };

  const handleEvent = (eventName: string, e: MessageEvent) => {
    if (eventName === "COMPETITOR_PUBLIC_INFO_UPDATED") {
      const event = competitorPublicInfoUpdatedEventSchema.parse(
        JSON.parse(e.data),
      );

      queueEventHandler((competitors: Map<number, ScoreboardEntry>) => {
        let competitor = competitors.get(event.competitorId);
        if (!competitor) {
          competitor = createEmptyEntry(event.competitorId);
        }

        competitor.categoryId = event.categoryId;
        competitor.name = event.name;
        competitor.withdrawnFromFinals = event.withdrawnFromFinals;
        competitor.disqualified = event.disqualified;
        competitor.scrubbedAt = event.scrubbedAt;

        competitors.set(event.competitorId, { ...competitor });
      });
    }

    if (eventName === "[]COMPETITOR_RESULT_UPDATED") {
      const events = z
        .array(competitorResultUpdatedEventSchema)
        .parse(JSON.parse(e.data));

      for (const event of events) {
        queueEventHandler((competitors: Map<number, ScoreboardEntry>) => {
          let competitor = competitors.get(event.competitorId);
          if (!competitor) {
            competitor = createEmptyEntry(event.competitorId);
          }

          const result: Result = {
            competitorId: event.competitorId,
            score: event.score,
            placement: event.placement,
            finalist: event.finalist,
            rankOrder: event.rankOrder,
            timestamp: event.timestamp,
          };

          competitor.result = result;
          competitors.set(event.competitorId, { ...competitor });
        });
      }
    }
  };

  const createEmptyEntry = (competitorId: number): ScoreboardEntry => ({
    competitorId: competitorId,
    categoryId: 0,
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
