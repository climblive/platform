# Carousel Item

**Full documentation:** https://webawesome.com/docs/components/carousel-item


`<wa-carousel-item>` Experimental Since 2.0

Carousel items represent individual slides within a carousel.

This component must be used as a child of `<wa-carousel>`. Please see the [Carousel docs](https://webawesome.com/docs/components/carousel) to see examples of this component in action.

## Importing

If you're using the autoloader or a hosted project, components load on demand — no manual import needed. To cherry-pick a component manually, use one of the following snippets.

\*\*CDN\*\*

Import this component directly from the CDN:

```js
import 'https://ka-f.webawesome.com/webawesome@3.7.0/components/carousel-item/carousel-item.js';
```

\*\*npm\*\*

After installing Web Awesome via npm, import this component:

```js
import '@awesome.me/webawesome/dist/components/carousel-item/carousel-item.js';
```

\*\*Self-Hosted\*\*

If you're self-hosting Web Awesome, import this component from your server:

```js
import './webawesome/dist/components/carousel-item/carousel-item.js';
```

\*\*React\*\*

To import this component for React 18 or below, use the following code:

```js
import WaCarouselItem from '@awesome.me/webawesome/dist/react/carousel-item/index.js';
```

## Slots

Learn more about [using slots](https://webawesome.com/docs/usage/#slots).

| Name | Description |
| --- | --- |
| (default) | The carousel item's content.. |

## Attributes & Properties

Learn more about [attributes and properties](https://webawesome.com/docs/usage/#attributes-and-properties).

| Name | Description | Reflects |
| --- | --- | --- |
| \`css\` | \`CSSResultGroup \\| undefined\` One or more CSSResultGroup to include in the component's shadow root. Host styles are automatically prepended. Type Default styles | | |

## CSS custom properties

Learn more about [CSS custom properties](https://webawesome.com/docs/usage/#custom-properties).

| Name | Description |
| --- | --- |
| \`--aspect-ratio\` | The slide's aspect ratio. Inherited from the carousel by default. |

**Need a hand?** Report a bug Ask for help