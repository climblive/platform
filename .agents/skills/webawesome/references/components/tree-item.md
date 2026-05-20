# Tree Item

**Full documentation:** https://webawesome.com/docs/components/tree-item


`<wa-tree-item>` Stable Since 2.0

Tree items represent a single hierarchical node inside a tree, and can contain nested items that expand and collapse.

This component must be used as a child of `<wa-tree>`. Please see the [Tree docs](https://webawesome.com/docs/components/tree) to see examples of this component in action.

## Importing

If you're using the autoloader or a hosted project, components load on demand — no manual import needed. To cherry-pick a component manually, use one of the following snippets.

\*\*CDN\*\*

Import this component directly from the CDN:

```js
import 'https://ka-f.webawesome.com/webawesome@3.7.0/components/tree-item/tree-item.js';
```

\*\*npm\*\*

After installing Web Awesome via npm, import this component:

```js
import '@awesome.me/webawesome/dist/components/tree-item/tree-item.js';
```

\*\*Self-Hosted\*\*

If you're self-hosting Web Awesome, import this component from your server:

```js
import './webawesome/dist/components/tree-item/tree-item.js';
```

\*\*React\*\*

To import this component for React 18 or below, use the following code:

```js
import WaTreeItem from '@awesome.me/webawesome/dist/react/tree-item/index.js';
```

## Slots

Learn more about [using slots](https://webawesome.com/docs/usage/#slots).

| Name | Description |
| --- | --- |
| (default) | The default slot. |
| \`collapse-icon\` | The icon to show when the tree item is collapsed. |
| \`expand-icon\` | The icon to show when the tree item is expanded. |

## Attributes & Properties

Learn more about [attributes and properties](https://webawesome.com/docs/usage/#attributes-and-properties).

| Name | Description | Reflects |
| --- | --- | --- |
| \`css\` | \`CSSResultGroup \\| undefined\` One or more CSSResultGroup to include in the component's shadow root. Host styles are automatically prepended. Type Default styles | | |
| \`disabled\` disabled | \`boolean\` Disables the tree item. Type Default false | | |
| \`expanded\` expanded | \`boolean\` Expands the tree item. Type Default false | | |
| \`lazy\` lazy | \`boolean\` Enables lazy loading behavior. Type Default false | | |
| \`selected\` selected | \`boolean\` Draws the tree item in a selected state. Type Default false | | |

## Methods

Learn more about [methods](https://webawesome.com/docs/usage/#methods).

| Name | Description | Arguments |
| --- | --- | --- |
| \`getChildrenItems()\` | Gets all the nested tree items in this node. | \`{ includeDisabled = true }: { includeDisabled?: boolean }\` |

## Events

Learn more about [events](https://webawesome.com/docs/usage/#events).

| Name | Description |
| --- | --- |
| \`wa-after-collapse\` | Emitted after the tree item collapses and all animations are complete. |
| \`wa-after-expand\` | Emitted after the tree item expands and all animations are complete. |
| \`wa-collapse\` | Emitted when the tree item collapses. |
| \`wa-expand\` | Emitted when the tree item expands. |
| \`wa-lazy-change\` | Emitted when the tree item's lazy state changes. |
| \`wa-lazy-load\` | \`lazy\` Emitted when a item is selected. Use this event to asynchronously load data and append items to the tree before expanding. After appending new items, remove the lazy attribute to remove the loading state and update the tree. |

## CSS custom properties

Learn more about [CSS custom properties](https://webawesome.com/docs/usage/#custom-properties).

| Name | Description |
| --- | --- |
| \`--hide-duration\` | \`200ms\` The animation duration when collapsing tree items. Default |
| \`--show-duration\` | \`200ms\` The animation duration when expanding tree items. Default |

## Custom States

Learn more about [custom states](https://webawesome.com/docs/usage/#custom-states).

| Name | Description | CSS selector |
| --- | --- | --- |
| \`disabled\` | Applied when the tree item is disabled. | \`:state(disabled)\` |
| \`expanded\` | Applied when the tree item is expanded. | \`:state(expanded)\` |
| \`indeterminate\` | Applied when the selection is indeterminate. | \`:state(indeterminate)\` |
| \`selected\` | Applied when the tree item is selected. | \`:state(selected)\` |

## CSS parts

Learn more about [CSS parts](https://webawesome.com/docs/usage/#css-parts).

| Name | Description | CSS selector |
| --- | --- | --- |
| \`base\` | The component's base wrapper. | \`::part(base)\` |
| \`checkbox\` | The checkbox that shows when using multiselect. | \`::part(checkbox)\` |
| \`checkbox\_\_base\` | \`base\` The checkbox's exported part. | \`::part(checkbox\_\_base)\` |
| \`checkbox\_\_checked-icon\` | \`checked-icon\` The checkbox's exported part. | \`::part(checkbox\_\_checked-icon)\` |
| \`checkbox\_\_control\` | \`control\` The checkbox's exported part. | \`::part(checkbox\_\_control)\` |
| \`checkbox\_\_indeterminate-icon\` | \`indeterminate-icon\` The checkbox's exported part. | \`::part(checkbox\_\_indeterminate-icon)\` |
| \`checkbox\_\_label\` | \`label\` The checkbox's exported part. | \`::part(checkbox\_\_label)\` |
| \`children\` | The container that wraps the tree item's nested children. | \`::part(children)\` |
| \`expand-button\` | The container that wraps the tree item's expand button and spinner. | \`::part(expand-button)\` |
| \`indentation\` | The tree item's indentation container. | \`::part(indentation)\` |
| \`item\` | The tree item's container. This element wraps everything except slotted tree item children. | \`::part(item)\` |
| \`label\` | The tree item's label. | \`::part(label)\` |
| \`spinner\` | The spinner that shows when a lazy tree item is in the loading state. | \`::part(spinner)\` |
| \`spinner\_\_base\` | The spinner's base part. | \`::part(spinner\_\_base)\` |

## Dependencies

This component automatically imports the following elements. Sub-dependencies, if any exist, will also be included in this list.

-   [`<wa-checkbox>`](https://webawesome.com/docs/components/checkbox)
-   [`<wa-icon>`](https://webawesome.com/docs/components/icon)
-   [`<wa-spinner>`](https://webawesome.com/docs/components/spinner)

**Need a hand?** Report a bug Ask for help