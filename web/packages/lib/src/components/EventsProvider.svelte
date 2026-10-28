<script lang="ts" module>
  import { createContext, type Snippet } from "svelte";
  import { SvelteMap } from "svelte/reactivity";
  import type { EventFilter } from "../models";
  import { base64UrlEncode, getApiUrl } from "../utils";

  const [get, set] = createContext<EventsClient>();

  type CallBack = (eventName: string, event: MessageEvent) => void;

  class EventsClient {
    #subscriptions: Map<symbol, [EventFilter[], CallBack]>;
    #filters: EventFilter[];

    constructor() {
      this.#subscriptions = new SvelteMap();
      this.#filters = $derived([
        ...this.#subscriptions.values().flatMap(([filters]) => filters),
      ]);

      $effect(() => {
        if (this.#filters.length === 0) {
          return;
        }

        const eventSource = new EventSource(
          `${getApiUrl()}/events?f=${base64UrlEncode(JSON.stringify(this.#filters))}`,
        );

        const uniqueEventTypes = new Set(
          this.#filters.flatMap(({ eventTypes }) =>
            eventTypes ? Object.keys(eventTypes) : [],
          ),
        );

        const listeners: [string, (e: MessageEvent) => void][] = [];

        for (const eventType of uniqueEventTypes.values()) {
          const listener = (event: MessageEvent) => {
            for (const [, cb] of this.#subscriptions.values()) {
              cb(event.type, event);
            }
          };

          eventSource.addEventListener(eventType, listener);

          listeners.push([eventType, listener]);
        }

        return () => {
          listeners.forEach(([eventName, listener]) =>
            eventSource.removeEventListener(eventName, listener),
          );

          eventSource.close();
        };
      });
    }

    public subscribe(filters: EventFilter[], callback: CallBack): symbol {
      const subscriptionId = Symbol();

      this.#subscriptions.set(subscriptionId, [filters, callback]);

      return subscriptionId;
    }

    public unsubscibe(subscriptionId: symbol): void {
      this.#subscriptions.delete(subscriptionId);
    }
  }

  export const subscribeEvents = (
    filters: EventFilter[],
    callback: CallBack,
  ) => {
    const client = get();
    return client.subscribe(filters, callback);
  };

  export const unsubscibeEvents = (subscriptionId: symbol) => {
    const client = get();
    return client.unsubscibe(subscriptionId);
  };
</script>

<script lang="ts">
  type Props = {
    children: Snippet;
  };

  const { children }: Props = $props();

  set(new EventsClient());
</script>

{@render children()}
