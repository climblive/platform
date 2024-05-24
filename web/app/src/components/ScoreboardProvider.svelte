<script lang="ts">
  import { Client } from "@stomp/stompjs";
  import { onDestroy, onMount, setContext } from "svelte";
  import { writable } from "svelte/store";
  import { ApiClient, configData } from "@climblive/shared";
  import type {
    ScoreboardContender,
    ScoreboardUpdate,
  } from "@climblive/shared/models";
  import type { RankedContender } from "@/types";

  export let contestId: number;
  export let numFinalists: number;

  let wsClient: Client;
  let initialized = false;

  const contenders: Map<number, ScoreboardContender> = new Map();
  const contendersToCompClass: Map<number, number> = new Map();
  const pendingUpdates: ScoreboardUpdate[] = [];

  const resultsStore = writable<Map<number, RankedContender[]>>(new Map());

  setContext("scoreboard", resultsStore);

  onMount(async () => {
    const scoreboard = await ApiClient.getInstance().getScoreboard(contestId);

    for (const { compClass, contenders: results } of scoreboard.scores) {
      for (const contender of results) {
        handleUpdate({ compClassId: compClass.id, contender });
      }
    }

    pendingUpdates.forEach(handleUpdate);

    initialized = true;
    $resultsStore = calculateResults();
  });

  const handleUpdate = ({ compClassId, contender }: ScoreboardUpdate) => {
    contenders.set(contender.contenderId, contender);
    contendersToCompClass.set(contender.contenderId, compClassId);
  };

  onMount(() => {
    wsClient = new Client({
      brokerURL: configData.WSS_URL,
      heartbeatIncoming: 4000,
      heartbeatOutgoing: 4000,
    });

    wsClient.activate();
    wsClient.onConnect = () => {
      wsClient.subscribe(
        `/topic/contest/${contestId}/scoreboard`,
        (message) => {
          const { compClassId, item } = JSON.parse(message.body) as {
            compClassId: number;
            item: ScoreboardContender;
          };

          if (!initialized) {
            pendingUpdates.push({ compClassId, contender: item });
          } else {
            handleUpdate({ compClassId, contender: item });
            $resultsStore = calculateResults();
          }
        }
      );
    };
  });

  onDestroy(() => {
    wsClient.deactivate();
  });

  const calculateResults = () => {
    const results: Map<number, RankedContender[]> = new Map();

    for (const contender of contenders.values()) {
      const compClassId = contendersToCompClass.get(contender.contenderId);
      if (!compClassId) {
        continue;
      }

      let classResults = results.get(compClassId);
      if (classResults === undefined) {
        classResults = [];
        results.set(compClassId, classResults);
      }

      classResults.push({
        ...contender,
        finalist: false,
        order: 0,
        placement: 0,
      });
    }

    for (const contenders of results.values()) {
      rankContenders(contenders, numFinalists);
    }

    return results;
  };

  const rankContenders = (
    contenders: RankedContender[],
    numFinalists: number
  ) => {
    const sortedContenders = contenders.toSorted(
      (c1, c2) => c2.qualifyingScore - c1.qualifyingScore
    );

    let index = 0;
    let placementCounter = 0;
    let prevScore = NaN;
    let gap = 1;
    let countedFinalists = 0;
    let highestFinalistPlacement = 0;

    for (const contender of sortedContenders) {
      contender.order = index++;

      if (prevScore !== contender.qualifyingScore) {
        placementCounter += gap;
        gap = 1;
      } else {
        gap += 1;
      }

      prevScore = contender.qualifyingScore;

      contender.placement = placementCounter;

      if (
        countedFinalists < numFinalists ||
        contender.placement === highestFinalistPlacement
      ) {
        contender.finalist = contender.qualifyingScore > 0;
        countedFinalists += 1;
        highestFinalistPlacement = placementCounter;
      }
    }
  };
</script>

<slot />
