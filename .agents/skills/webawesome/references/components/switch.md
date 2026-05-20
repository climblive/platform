# Switch

**Full documentation:** https://webawesome.com/docs/components/switch


`<wa-switch>` Stable Since 2.0

Switches toggle a single setting on or off and apply the change immediately, without requiring a form submission.

```html
<wa-switch>Switch</wa-switch>
```

This component works with standard `<form>` elements. Please refer to the section on [form controls](https://webawesome.com/docs/form-controls) to learn more about form submission and client-side validation.

## Examples

### Checked

Use the `checked` attribute to activate the switch.

```html
<wa-switch checked>Checked</wa-switch>
```

The `checked` attribute is the initial value and does not reflect changes, consistent with native checkboxes. To toggle the checked state with JavaScript, use the `checked` property instead. To target checked switches with CSS, use the `:state(checked)` selector.

### Disabled

Use the `disabled` attribute to disable the switch.

```html
<wa-switch disabled>Disabled</wa-switch>
```

### Sizes

Use the `size` attribute to change a switch's size.

```html
<wa-switch size="xs">Extra Small</wa-switch>
<br />
<wa-switch size="s">Small</wa-switch>
<br />
<wa-switch size="m">Medium</wa-switch>
<br />
<wa-switch size="l">Large</wa-switch>
<br />
<wa-switch size="xl">Extra Large</wa-switch>
```

### Hint

Add descriptive hint to a switch with the `hint` attribute. For hints that contain HTML, use the `hint` slot instead.

```html
<wa-switch hint="What should the user know about the switch?">Label</wa-switch>
```

### Custom Styles

Use the available custom properties to change how the switch is styled.

```html
<wa-switch style="--width: 80px; --height: 40px; --thumb-size: 36px;">Really big</wa-switch>
```

## Importing

If you're using the autoloader or a hosted project, components load on demand — no manual import needed. To cherry-pick a component manually, use one of the following snippets.

\*\*CDN\*\*

Import this component directly from the CDN:

```js
import 'https://ka-f.webawesome.com/webawesome@3.7.0/components/switch/switch.js';
```

\*\*npm\*\*

After installing Web Awesome via npm, import this component:

```js
import '@awesome.me/webawesome/dist/components/switch/switch.js';
```

\*\*Self-Hosted\*\*

If you're self-hosting Web Awesome, import this component from your server:

```js
import './webawesome/dist/components/switch/switch.js';
```

\*\*React\*\*

To import this component for React 18 or below, use the following code:

```js
import WaSwitch from '@awesome.me/webawesome/dist/react/switch/index.js';
```

## Slots

Learn more about [using slots](https://webawesome.com/docs/usage/#slots).

| Name | Description |
| --- | --- |
| (default) | The switch's label. |
| \`hint\` | \`hint\` Text that describes how to use the switch. Alternatively, you can use the attribute. |

## Attributes & Properties

Learn more about [attributes and properties](https://webawesome.com/docs/usage/#attributes-and-properties).

| Name | Description | Reflects |
| --- | --- | --- |
| \`checked\` | Draws the checkbox in a checked state. | | |
| \`css\` | \`CSSResultGroup \\| undefined\` One or more CSSResultGroup to include in the component's shadow root. Host styles are automatically prepended. Type Default \[formControlStyles, sizeStyles, styles\] | | |
| \`defaultChecked\` checked | \`boolean\` The default value of the form control. Primarily used for resetting the form control. Type | | |
| \`disabled\` disabled | \`boolean\` Disables the switch. Type Default false | | |
| \`form\` | \`

\` By default, form controls are associated with the nearest containing element. This attribute allows you to place the form control outside of a form and associate it with the form that has this id. The form must be in the same document or shadow root for this to work. Type HTMLFormElement \\| null | | |
| \`hint\` hint | \`hint\` The switch's . If you need to display HTML, use the hint slot instead. Type string Default '' | | |
| \`name\` name | \`string \\| null\` The name of the switch, submitted as a name/value pair with form data. Type Default null | | |
| \`required\` required | \`boolean\` Makes the switch a required field. Type Default false | | |
| \`size\` size | \`'xs' \\| 's' \\| 'm' \\| 'l' \\| 'xl' \\| 'small' \\| 'medium' \\| 'large'\` The switch's size. Type Default 'm' | | |
| \`validationTarget\` | \`undefined \\| HTMLElement\` Override this to change where constraint validation popups are anchored. Type | | |
| \`validators\` | \`observedAttributes\` Validators are static because they have , essentially attributes to "watch" for changes. Whenever these attributes change, we want to be notified and update the validator. Type Validator\[\] Default \[\] | | |
| \`value\` value | \`string \\| null\` The value of the switch, submitted as a name/value pair with form data. Type | | |
| \`withHint\` with-hint | \`true\` Only required for SSR. Set to if you're slotting in a hint element so the server-rendered markup includes the hint before the component hydrates on the client. Type boolean Default false | | |

## Methods

Learn more about [methods](https://webawesome.com/docs/usage/#methods).

| Name | Description | Arguments |
| --- | --- | --- |
| \`blur()\` | Removes focus from the switch. | |
| \`click()\` | Simulates a click on the switch. | |
| \`focus()\` | Sets focus on the switch. | \`options: FocusOptions\` |
| \`formStateRestoreCallback()\` | Called when the browser is trying to restore element’s state to state in which case reason is "restore", or when the browser is trying to fulfill autofill on behalf of user in which case reason is "autocomplete". In the case of "restore", state is a string, File, or FormData object previously set as the second argument to setFormValue. | \`state: string \\| File \\| FormData \\| null, reason: 'autocomplete' \\| 'restore'\` |
| \`resetValidity()\` | Reset validity is a way of removing manual custom errors and native validation. | |
| \`setCustomValidity()\` | Do not use this when creating a "Validator". This is intended for end users of components. We track manually defined custom errors so we don't clear them on accident in our validators. | \`message: string\` |

## Events

Learn more about [events](https://webawesome.com/docs/usage/#events).

| Name | Description |
| --- | --- |
| \`blur\` | Emitted when the control loses focus. |
| \`change\` | Emitted when the control's checked state changes. |
| \`focus\` | Emitted when the control gains focus. |
| \`input\` | Emitted when the control receives input. |
| \`wa-invalid\` | Emitted when the form control has been checked for validity and its constraints aren't satisfied. |

## CSS custom properties

Learn more about [CSS custom properties](https://webawesome.com/docs/usage/#custom-properties).

| Name | Description |
| --- | --- |
| \`--height\` | The height of the switch. |
| \`--thumb-size\` | The size of the thumb. |
| \`--width\` | The width of the switch. |

## CSS parts

Learn more about [CSS parts](https://webawesome.com/docs/usage/#css-parts).

| Name | Description | CSS selector |
| --- | --- | --- |
| \`base\` | The component's base wrapper. | \`::part(base)\` |
| \`control\` | The control that houses the switch's thumb. | \`::part(control)\` |
| \`hint\` | The hint's wrapper. | \`::part(hint)\` |
| \`label\` | The switch's label. | \`::part(label)\` |
| \`thumb\` | The switch's thumb. | \`::part(thumb)\` |

**Need a hand?** Report a bug Ask for help