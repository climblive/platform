# Zoomable Frame

**Full documentation:** https://webawesome.com/docs/components/zoomable-frame


`<wa-zoomable-frame>` Stable Since 3.0

Zoomable frames embed iframe content with built-in controls for zooming, panning, and managing interaction.

```html
<wa-zoomable-frame src="/examples/themes/showcase" zoom="0.5"> </wa-zoomable-frame>
```

## Examples

### Loading external content

Use the `src` attribute to embed external websites or resources. The URL must be accessible, and cross-origin restrictions may apply due to the Same-Origin Policy, potentially limiting access to the iframe's content.

```html
<wa-zoomable-frame src="https://example.com/"> </wa-zoomable-frame>
```

The zoomable frame fills 100% width by default with a 16:9 aspect ratio. Customize this using the `aspect-ratio` CSS property.

```html
<wa-zoomable-frame src="https://example.com/" style="aspect-ratio: 4/3;"> </wa-zoomable-frame>
```

Use the `srcdoc` attribute or property to display custom HTML content directly within the iframe, perfect for rendering inline content without external resources.

```html
<wa-zoomable-frame srcdoc="<html><body><h1>Hello, World!</h1><p>This is inline content.</p></body></html>">
</wa-zoomable-frame>
```

When both `src` and `srcdoc` are specified, `srcdoc` takes precedence.

### Controlling zoom behavior

Set the `zoom` attribute to control the frame's zoom level. Use `1` for 100%, `2` for 200%, `0.5` for 50%, and so on.

Define specific zoom increments with the `zoom-levels` attribute using space-separated percentages and decimal values like `zoom-levels="0.25 0.5 75% 100%"`.

```html
<wa-zoomable-frame src="/examples/themes/showcase" zoom="0.5" zoom-levels="50% 0.75 100%"> </wa-zoomable-frame>
```

### Hiding zoom controls

Add the `without-controls` attribute to hide the zoom control interface from the frame.

```html
<wa-zoomable-frame src="/examples/themes/showcase" without-controls zoom="0.5"> </wa-zoomable-frame>
```

### Preventing user interaction

Apply the `without-interaction` attribute to make the frame non-interactive. Note that this prevents keyboard navigation into the frame, which may impact accessibility for some users.

```html
<wa-zoomable-frame src="/examples/themes/showcase" zoom="0.5" without-interaction> </wa-zoomable-frame>
```

### Enabling theme sync

By default, the frame does not sync theme classes into the iframe. Add the `with-theme-sync` attribute to mirror the host page's light/dark mode and theme selector classes (such as `wa-theme-*`, `wa-brand-*`, and `wa-palette-*`) into the iframe document. This is useful when the iframe renders Web Awesome styles that should match the host page's theme.

```html
<wa-zoomable-frame src="/examples/themes/showcase" zoom="0.5" with-theme-sync> </wa-zoomable-frame>
```

## Importing

If you're using the autoloader or a hosted project, components load on demand — no manual import needed. To cherry-pick a component manually, use one of the following snippets.

\*\*CDN\*\*

Import this component directly from the CDN:

```js
import 'https://ka-f.webawesome.com/webawesome@3.7.0/components/zoomable-frame/zoomable-frame.js';
```

\*\*npm\*\*

After installing Web Awesome via npm, import this component:

```js
import '@awesome.me/webawesome/dist/components/zoomable-frame/zoomable-frame.js';
```

\*\*Self-Hosted\*\*

If you're self-hosting Web Awesome, import this component from your server:

```js
import './webawesome/dist/components/zoomable-frame/zoomable-frame.js';
```

\*\*React\*\*

To import this component for React 18 or below, use the following code:

```js
import WaZoomableFrame from '@awesome.me/webawesome/dist/react/zoomable-frame/index.js';
```

## Slots

Learn more about [using slots](https://webawesome.com/docs/usage/#slots).

| Name | Description |
| --- | --- |
| \`zoom-in-icon\` | The slot that contains the zoom in icon. |
| \`zoom-out-icon\` | The slot that contains the zoom out icon. |

## Attributes & Properties

Learn more about [attributes and properties](https://webawesome.com/docs/usage/#attributes-and-properties).

| Name | Description | Reflects |
| --- | --- | --- |
| \`allowfullscreen\` allowfullscreen | \`boolean\` Allows fullscreen mode. Type Default false | | |
| \`contentDocument\` | \`document\` Returns the internal iframe's object. (Readonly property) Type Document \\| null | | |
| \`contentWindow\` | \`window\` Returns the internal iframe's object. (Readonly property) Type Window \\| null | | |
| \`css\` | \`CSSResultGroup \\| undefined\` One or more CSSResultGroup to include in the component's shadow root. Host styles are automatically prepended. Type Default styles | | |
| \`loading\` loading | \`'eager' \\| 'lazy'\` Controls iframe loading behavior. Type Default 'eager' | | |
| \`referrerpolicy\` referrerpolicy | \`string\` Controls referrer information. Type | | |
| \`sandbox\` sandbox | \`string\` Security restrictions for the iframe. Type | | |
| \`src\` src | \`string\` The URL of the content to display. Type | | |
| \`srcdoc\` srcdoc | \`string\` Inline HTML to display. Type | | |
| \`withoutControls\` without-controls | \`boolean\` Removes the zoom controls. Type Default false | | |
| \`withoutInteraction\` without-interaction | \`boolean\` Disables interaction when present. Type Default false | | |
| \`withThemeSync\` with-theme-sync | \`boolean\` Enables automatic theme syncing (light/dark mode and theme selector classes) from the host document to the iframe. Type Default false | | |
| \`zoom\` zoom | \`number\` The current zoom of the frame, e.g. 0 = 0% and 1 = 100%. Type Default 1 | | |
| \`zoomLevels\` zoom-levels | \`string\` The zoom levels to step through when using zoom controls. This does not restrict programmatic changes to the zoom. Type Default '25% 50% 75% 100% 125% 150% 175% 200%' | | |

## Methods

Learn more about [methods](https://webawesome.com/docs/usage/#methods).

| Name | Description | Arguments |
| --- | --- | --- |
| \`zoomIn()\` | Zooms in to the next available zoom level. | |
| \`zoomOut()\` | Zooms out to the previous available zoom level. | |

## Events

Learn more about [events](https://webawesome.com/docs/usage/#events).

| Name | Description |
| --- | --- |
| \`error\` | Emitted from the internal iframe when it fails to load. |
| \`load\` | Emitted when the internal iframe when it finishes loading. |

## CSS parts

Learn more about [CSS parts](https://webawesome.com/docs/usage/#css-parts).

| Name | Description | CSS selector |
| --- | --- | --- |
| \`controls\` | The container that surrounds zoom control buttons. | \`::part(controls)\` |
| \`iframe\` | \`\` The internal element. | \`::part(iframe)\` |
| \`zoom-in-button\` | The zoom in button. | \`::part(zoom-in-button)\` |
| \`zoom-out-button\` | The zoom out button. | \`::part(zoom-out-button)\` |

## Dependencies

This component automatically imports the following elements. Sub-dependencies, if any exist, will also be included in this list.

-   [`<wa-icon>`](https://webawesome.com/docs/components/icon)

**Need a hand?** Report a bug Ask for help