# Breadcrumb Item

**Full documentation:** https://webawesome.com/docs/components/breadcrumb-item


`<wa-breadcrumb-item>` Stable Since 2.0

Breadcrumb items represent individual links inside a breadcrumb, typically one per level of the site hierarchy.

This component must be used as a child of `<wa-breadcrumb>`. Please see the [Breadcrumb docs](https://webawesome.com/docs/components/breadcrumb) to see examples of this component in action.

## Importing

If you're using the autoloader or a hosted project, components load on demand — no manual import needed. To cherry-pick a component manually, use one of the following snippets.

\*\*CDN\*\*

Import this component directly from the CDN:

```js
import 'https://ka-f.webawesome.com/webawesome@3.7.0/components/breadcrumb-item/breadcrumb-item.js';
```

\*\*npm\*\*

After installing Web Awesome via npm, import this component:

```js
import '@awesome.me/webawesome/dist/components/breadcrumb-item/breadcrumb-item.js';
```

\*\*Self-Hosted\*\*

If you're self-hosting Web Awesome, import this component from your server:

```js
import './webawesome/dist/components/breadcrumb-item/breadcrumb-item.js';
```

\*\*React\*\*

To import this component for React 18 or below, use the following code:

```js
import WaBreadcrumbItem from '@awesome.me/webawesome/dist/react/breadcrumb-item/index.js';
```

## Slots

Learn more about [using slots](https://webawesome.com/docs/usage/#slots).

| Name | Description |
| --- | --- |
| (default) | The breadcrumb item's label. |
| \`end\` | \`

## Attributes & Properties

Learn more about [attributes and properties](https://webawesome.com/docs/usage/#attributes-and-properties).

| Name | Description | Reflects |
| --- | --- | --- |
| \`css\` | \`CSSResultGroup \\| undefined\` One or more CSSResultGroup to include in the component's shadow root. Host styles are automatically prepended. Type Default styles | | |
| \`href\` href | \`string \\| undefined\` Optional URL to direct the user to when the breadcrumb item is activated. When set, a link will be rendered internally. When unset, a button will be rendered instead. Type | | |
| \`rel\` rel | \`rel\` The attribute to use on the link. Only used when href is set. Type string Default 'noreferrer noopener' | | |
| \`target\` target | \`href\` Tells the browser where to open the link. Only used when is set. Type '\_blank' \\| '\_parent' \\| '\_self' \\| '\_top' \\| undefined | | |

## CSS parts

Learn more about [CSS parts](https://webawesome.com/docs/usage/#css-parts).

| Name | Description | CSS selector |
| --- | --- | --- |
| \`end\` | \`end\` The container that wraps the slot. | \`::part(end)\` |
| \`label\` | The breadcrumb item's label. | \`::part(label)\` |
| \`separator\` | The container that wraps the separator. | \`::part(separator)\` |
| \`start\` | \`start\` The container that wraps the slot. | \`::part(start)\` |

**Need a hand?** Report a bug Ask for help