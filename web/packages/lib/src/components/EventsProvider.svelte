<script lang="ts" module>
  import { createContext, type Snippet } from "svelte";
  import { SvelteMap } from "svelte/reactivity";
  import type { EventFilter } from "../models";
  import { base64UrlEncode } from "../utils";

  const [get, set] = createContext<EventsClient>();

  type CallBack = (eventName: string, event: MessageEvent) => void;
  type ConnectionHandlers = {
    onOpen?: () => void;
    onError?: () => void;
  };

  const invokeCallback = (callback: (() => void) | undefined) => {
    try {
      callback?.();
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error("Event subscriber callback failed", error);
    }
  };

  class EventsClient {
    #subscriptions: Map<symbol, [EventFilter[], CallBack, ConnectionHandlers]>;
    #filters: EventFilter[];

    constructor() {
      this.#subscriptions = new SvelteMap();
      this.#filters = $derived(
        [...this.#subscriptions.values()].flatMap(([filters]) => filters),
      );

      $effect(() => {
        if (this.#filters.length === 0) {
          return;
        }

        const listeners: [string, (e: MessageEvent) => void][] = [];
        let eventSource: EventSource;
        let retryTimeout: ReturnType<typeof setTimeout> | undefined;

        const connect = () => {
          eventSource = new EventSource(
            `/api/events?f=${base64UrlEncode(JSON.stringify(this.#filters))}`,
          );

          eventSource.onopen = () => {
            for (const [, , handlers] of this.#subscriptions.values()) {
              invokeCallback(handlers.onOpen);
            }
          };

          eventSource.onerror = () => {
            for (const [, , handlers] of this.#subscriptions.values()) {
              invokeCallback(handlers.onError);
            }

            if (eventSource.readyState === EventSource.CLOSED) {
              disconnect();
              retryTimeout = setTimeout(connect, 5000);
            }
          };

          for (const [eventType, listener] of listeners) {
            eventSource.addEventListener(eventType, listener);
          }
        };

        const disconnect = () => {
          eventSource.onopen = null;
          eventSource.onerror = null;
          for (const [eventType, listener] of listeners) {
            eventSource.removeEventListener(eventType, listener);
          }
          eventSource.close();
        };

        for (const [id, [filters, callback]] of this.#subscriptions) {
          const eventTypes = new Set(
            filters.flatMap(({ eventTypes }) => Object.keys(eventTypes ?? {})),
          );

          for (const eventType of eventTypes) {
            const listener = (event: MessageEvent) => {
              if (this.#subscriptions.has(id)) {
                invokeCallback(() => callback(event.type, event));
              }
            };

            listeners.push([eventType, listener]);
          }
        }

        connect();

        return () => {
          clearTimeout(retryTimeout);
          disconnect();
        };
      });
    }

    public subscribe(
      filters: EventFilter[],
      callback: CallBack,
      handlers: ConnectionHandlers = {},
    ): symbol {
      const subscriptionId = Symbol();

      this.#subscriptions.set(subscriptionId, [filters, callback, handlers]);

      return subscriptionId;
    }

    public unsubscribe(subscriptionId: symbol): void {
      this.#subscriptions.delete(subscriptionId);
    }
  }

  export const useEvents = () => {
    const client = get();
    return {
      subscribeEvents: client.subscribe.bind(client),
      unsubscribeEvents: client.unsubscribe.bind(client),
    };
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
