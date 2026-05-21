# Tab

**Full documentation:** https://webawesome.com/docs/components/tab


`<wa-tab>` Stable Since 2.0

Tabs label and activate an individual panel inside a tab group.

This component must be used as a child of `<wa-tab-group>`. Please see the [Tab Group docs](https://webawesome.com/docs/components/tab-group) to see examples of this component in action.

## Importing

If you're using the autoloader or a hosted project, components load on demand — no manual import needed. To cherry-pick a component manually, use one of the following snippets.

\*\*CDN\*\*

Import this component directly from the CDN:

```js
import 'https://ka-f.webawesome.com/webawesome@3.7.0/components/tab/tab.js';
```

\*\*npm\*\*

After installing Web Awesome via npm, import this component:

```js
import '@awesome.me/webawesome/dist/components/tab/tab.js';
```

\*\*Self-Hosted\*\*

If you're self-hosting Web Awesome, import this component from your server:

```js
import './webawesome/dist/components/tab/tab.js';
```

\*\*React\*\*

To import this component for React 18 or below, use the following code:

```js
import WaTab from '@awesome.me/webawesome/dist/react/tab/index.js';
```

## Slots

Learn more about [using slots](https://webawesome.com/docs/usage/#slots).

| Name | Description |
| --- | --- |
| (default) | The tab's label. |

## Attributes & Properties

Learn more about [attributes and properties](https://webawesome.com/docs/usage/#attributes-and-properties).

| Name | Description | Reflects |
| --- | --- | --- |
| \`css\` | \`CSSResultGroup \\| undefined\` One or more CSSResultGroup to include in the component's shadow root. Host styles are automatically prepended. Type Default styles | | |
| \`disabled\` disabled | \`boolean\` Disables the tab and prevents selection. Type Default false | | |
| \`panel\` panel | \`string\` The name of the tab panel this tab is associated with. The panel must be located in the same tab group. Type Default '' | | |

## CSS parts

Learn more about [CSS parts](https://webawesome.com/docs/usage/#css-parts).

| Name | Description | CSS selector |
| --- | --- | --- |
| \`base\` | The component's base wrapper. | \`::part(base)\` |

**Need a hand?** Report a bug Ask for help