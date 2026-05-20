# Tab Panel

**Full documentation:** https://webawesome.com/docs/components/tab-panel


`<wa-tab-panel>` Stable Since 2.0

Tab panels hold the content shown for a single tab inside a tab group.

This component must be used as a child of `<wa-tab-group>`. Please see the [Tab Group docs](https://webawesome.com/docs/components/tab-group) to see examples of this component in action.

## Importing

If you're using the autoloader or a hosted project, components load on demand — no manual import needed. To cherry-pick a component manually, use one of the following snippets.

\*\*CDN\*\*

Import this component directly from the CDN:

```js
import 'https://ka-f.webawesome.com/webawesome@3.7.0/components/tab-panel/tab-panel.js';
```

\*\*npm\*\*

After installing Web Awesome via npm, import this component:

```js
import '@awesome.me/webawesome/dist/components/tab-panel/tab-panel.js';
```

\*\*Self-Hosted\*\*

If you're self-hosting Web Awesome, import this component from your server:

```js
import './webawesome/dist/components/tab-panel/tab-panel.js';
```

\*\*React\*\*

To import this component for React 18 or below, use the following code:

```js
import WaTabPanel from '@awesome.me/webawesome/dist/react/tab-panel/index.js';
```

## Slots

Learn more about [using slots](https://webawesome.com/docs/usage/#slots).

| Name | Description |
| --- | --- |
| (default) | The tab panel's content. |

## Attributes & Properties

Learn more about [attributes and properties](https://webawesome.com/docs/usage/#attributes-and-properties).

| Name | Description | Reflects |
| --- | --- | --- |
| \`active\` active | \`boolean\` When true, the tab panel will be shown. Type Default false | | |
| \`css\` | \`CSSResultGroup \\| undefined\` One or more CSSResultGroup to include in the component's shadow root. Host styles are automatically prepended. Type Default styles | | |
| \`name\` name | \`string\` The tab panel's name. Type Default '' | | |

## CSS custom properties

Learn more about [CSS custom properties](https://webawesome.com/docs/usage/#custom-properties).

| Name | Description |
| --- | --- |
| \`--padding\` | The tab panel's padding. |

## CSS parts

Learn more about [CSS parts](https://webawesome.com/docs/usage/#css-parts).

| Name | Description | CSS selector |
| --- | --- | --- |
| \`base\` | The component's base wrapper. | \`::part(base)\` |

**Need a hand?** Report a bug Ask for help