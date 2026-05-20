# Dropdown Item

**Full documentation:** https://webawesome.com/docs/components/dropdown-item


`<wa-dropdown-item>` Stable Since 3.0

Dropdown items represent selectable entries within a dropdown menu, including standard actions, checkable items, and submenu triggers.

This component must be used as a child of `<wa-dropdown>`. Please see the [Dropdown docs](https://webawesome.com/docs/components/dropdown) to see examples of this component in action.

## Importing

If you're using the autoloader or a hosted project, components load on demand — no manual import needed. To cherry-pick a component manually, use one of the following snippets.

\*\*CDN\*\*

Import this component directly from the CDN:

```js
import 'https://ka-f.webawesome.com/webawesome@3.7.0/components/dropdown-item/dropdown-item.js';
```

\*\*npm\*\*

After installing Web Awesome via npm, import this component:

```js
import '@awesome.me/webawesome/dist/components/dropdown-item/dropdown-item.js';
```

\*\*Self-Hosted\*\*

If you're self-hosting Web Awesome, import this component from your server:

```js
import './webawesome/dist/components/dropdown-item/dropdown-item.js';
```

\*\*React\*\*

To import this component for React 18 or below, use the following code:

```js
import WaDropdownItem from '@awesome.me/webawesome/dist/react/dropdown-item/index.js';
```

## Slots

Learn more about [using slots](https://webawesome.com/docs/usage/#slots).

| Name | Description |
| --- | --- |
| (default) | The dropdown item's label. |
| \`details\` | Additional content or details to display after the label. |
| \`icon\` | An optional icon to display before the label. |
| \`submenu\` | \`\` Submenu items, typically elements, to create a nested menu. |

## Attributes & Properties

Learn more about [attributes and properties](https://webawesome.com/docs/usage/#attributes-and-properties).

| Name | Description | Reflects |
| --- | --- | --- |
| \`checked\` checked | \`type\` Set to true to check the dropdown item. Only valid when is checkbox. Type boolean Default false | | |
| \`css\` | \`CSSResultGroup \\| undefined\` One or more CSSResultGroup to include in the component's shadow root. Host styles are automatically prepended. Type Default styles | | |
| \`disabled\` disabled | \`boolean\` Disables the dropdown item. Type Default false | | |
| \`submenuOpen\` submenuOpen | \`boolean\` Whether the submenu is currently open. Type Default false | | |
| \`type\` type | \`checkbox\` Set to to make the item a checkbox. Type 'normal' \\| 'checkbox' Default 'normal' | | |
| \`value\` value | \`wa-select\` An optional value for the menu item. This is useful for determining which item was selected when listening to the dropdown's event. Type string | | |
| \`variant\` variant | \`'danger' \\| 'default'\` The type of menu item to render. Type Default 'default' | | |

## Methods

Learn more about [methods](https://webawesome.com/docs/usage/#methods).

| Name | Description | Arguments |
| --- | --- | --- |
| \`closeSubmenu()\` | Closes the submenu. | |
| \`openSubmenu()\` | Opens the submenu. | |

## Events

Learn more about [events](https://webawesome.com/docs/usage/#events).

| Name | Description |
| --- | --- |
| \`blur\` | Emitted when the dropdown item loses focus. |
| \`focus\` | Emitted when the dropdown item gains focus. |

## CSS parts

Learn more about [CSS parts](https://webawesome.com/docs/usage/#css-parts).

| Name | Description | CSS selector |
| --- | --- | --- |
| \`checkmark\` | \`\` The submenu indicator icon (a element). | \`::part(submenu-icon)\` |

## Dependencies

This component automatically imports the following elements. Sub-dependencies, if any exist, will also be included in this list.

-   [`<wa-icon>`](https://webawesome.com/docs/components/icon)

**Need a hand?** Report a bug Ask for help