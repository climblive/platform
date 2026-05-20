# Resize Observer

**Full documentation:** https://webawesome.com/docs/components/resize-observer


`<wa-resize-observer>` Stable Since 2.0

Resize observers watch their slotted elements for size changes and emit an event when they occur. Provides a thin, declarative interface to the browser's ResizeObserver API.

The resize observer will report changes to the dimensions of the elements it wraps through the `wa-resize` event. When emitted, a collection of [`ResizeObserverEntry`](https://developer.mozilla.org/en-US/docs/Web/API/ResizeObserverEntry) objects will be attached to `event.detail` that contains the target element and information about its dimensions.

```html
<div class="resize-observer-overview">
  <wa-resize-observer>
    <div>Resize this box and watch the console 👉</div>
  </wa-resize-observer>
</div>

<script>
  const container = document.querySelector('.resize-observer-overview');
  const resizeObserver = container.querySelector('wa-resize-observer');

  resizeObserver.addEventListener('wa-resize', event => {
    console.log(event.detail);
  });
</script>

<style>
  .resize-observer-overview div {
    display: flex;
    border: solid 2px var(--wa-color-surface-border);
    align-items: center;
    justify-content: center;
    text-align: center;
    padding: 4rem 2rem;
  }
</style>
```

## Importing

If you're using the autoloader or a hosted project, components load on demand — no manual import needed. To cherry-pick a component manually, use one of the following snippets.

\*\*CDN\*\*

Import this component directly from the CDN:

```js
import 'https://ka-f.webawesome.com/webawesome@3.7.0/components/resize-observer/resize-observer.js';
```

\*\*npm\*\*

After installing Web Awesome via npm, import this component:

```js
import '@awesome.me/webawesome/dist/components/resize-observer/resize-observer.js';
```

\*\*Self-Hosted\*\*

If you're self-hosting Web Awesome, import this component from your server:

```js
import './webawesome/dist/components/resize-observer/resize-observer.js';
```

\*\*React\*\*

To import this component for React 18 or below, use the following code:

```js
import WaResizeObserver from '@awesome.me/webawesome/dist/react/resize-observer/index.js';
```

## Slots

Learn more about [using slots](https://webawesome.com/docs/usage/#slots).

| Name | Description |
| --- | --- |
| (default) | One or more elements to watch for resizing. |

## Attributes & Properties

Learn more about [attributes and properties](https://webawesome.com/docs/usage/#attributes-and-properties).

| Name | Description | Reflects |
| --- | --- | --- |
| \`css\` | \`CSSResultGroup \\| undefined\` One or more CSSResultGroup to include in the component's shadow root. Host styles are automatically prepended. Type Default styles | | |
| \`disabled\` disabled | \`boolean\` Disables the observer. Type Default false | | |

## Events

Learn more about [events](https://webawesome.com/docs/usage/#events).

| Name | Description |
| --- | --- |
| \`wa-resize\` | Emitted when the element is resized. |

**Need a hand?** Report a bug Ask for help