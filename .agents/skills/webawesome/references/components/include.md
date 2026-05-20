# Include

**Full documentation:** https://webawesome.com/docs/components/include


`<wa-include>` Stable Since 2.0

Fetches an external HTML file and embeds its contents inline on the page. Useful for reusing shared markup like headers, footers, and partials across multiple pages.

Included files are asynchronously requested using `window.fetch()`. Requests are cached, so the same file can be included multiple times, but only one request will be made.

The included content will be inserted into the `<wa-include>` element's default slot so it can be easily accessed and styled through the light DOM.

```html
<wa-include src="https://shoelace.style/assets/examples/include.html"></wa-include>
```

## Examples

### Listening for Events

When an include file loads successfully, the `wa-load` event will be emitted. You can listen for this event to add custom loading logic to your includes.

If the request fails, the `wa-include-error` event will be emitted. In this case, `event.detail.status` will contain the resulting HTTP status code of the request, e.g. 404 (not found).

```html
<wa-include src="https://shoelace.style/assets/examples/include.html"></wa-include>

<script>
  const include = document.querySelector('wa-include');

  include.addEventListener('wa-load', event => {
    if (event.eventPhase === Event.AT_TARGET) {
      console.log('Success');
    }
  });

  include.addEventListener('wa-include-error', event => {
    if (event.eventPhase === Event.AT_TARGET) {
      console.log('Error', event.detail.status);
    }
  });
</script>
```

## Importing

If you're using the autoloader or a hosted project, components load on demand — no manual import needed. To cherry-pick a component manually, use one of the following snippets.

\*\*CDN\*\*

Import this component directly from the CDN:

```js
import 'https://ka-f.webawesome.com/webawesome@3.7.0/components/include/include.js';
```

\*\*npm\*\*

After installing Web Awesome via npm, import this component:

```js
import '@awesome.me/webawesome/dist/components/include/include.js';
```

\*\*Self-Hosted\*\*

If you're self-hosting Web Awesome, import this component from your server:

```js
import './webawesome/dist/components/include/include.js';
```

\*\*React\*\*

To import this component for React 18 or below, use the following code:

```js
import WaInclude from '@awesome.me/webawesome/dist/react/include/index.js';
```

## Attributes & Properties

Learn more about [attributes and properties](https://webawesome.com/docs/usage/#attributes-and-properties).

| Name | Description | Reflects |
| --- | --- | --- |
| \`allowScripts\` allow-scripts | \`boolean\` Allows included scripts to be executed. Be sure you trust the content you are including as it will be executed as code and can result in XSS attacks. Type Default false | | |
| \`css\` | \`CSSResultGroup \\| undefined\` One or more CSSResultGroup to include in the component's shadow root. Host styles are automatically prepended. Type Default styles | | |
| \`mode\` mode | \`'cors' \\| 'no-cors' \\| 'same-origin'\` The fetch mode to use. Type Default 'cors' | | |
| \`src\` src | \`string\` The location of the HTML file to include. Be sure you trust the content you are including as it will be executed as code and can result in XSS attacks. Type | | |

## Events

Learn more about [events](https://webawesome.com/docs/usage/#events).

| Name | Description |
| --- | --- |
| \`wa-include-error\` | Emitted when the included file fails to load due to an error. |
| \`wa-load\` | Emitted when the included file is loaded. |

**Need a hand?** Report a bug Ask for help