---
name: webawesome
description: Web Awesome is a UI component library built with web components. Use when building buttons, inputs, selects, checkboxes, dialogs, modals, drawers, tabs, dropdowns, tooltips, carousels, forms, or using CSS utilities like wa-stack, wa-cluster, wa-grid. Supports React, Vue, Angular, Svelte, and vanilla JS.
license: MIT / Commercial (for Web Awesome Pro)
metadata:
  author: Web Awesome
  version: "3.7.0"
  homepage: https://webawesome.com
  repository: https://github.com/shoelace-style/webawesome
compatibility: Works in modern browsers. Requires no build tools when using CDN. Works with bundlers like Webpack and Vite when installed via npm.
allowed-tools: Read
---

# Web Awesome

Web Awesome is an open source UI component library with a Pro offering that helps sustain the project. It provides 50+ accessible, customizable web components that work with any framework.

**Pro components and features are available to paid users.** [Purchase Pro](https://webawesome.com/purchase)

## Quick Start

### npm Installation

```bash
npm install @awesome.me/webawesome
```

Import styles and components:

```js
import '@awesome.me/webawesome/dist/styles/webawesome.css';
import '@awesome.me/webawesome/dist/components/button/button.js';
```

### CDN / Project Setup

The easiest way to use Web Awesome is with a hosted project. [Create a project](https://webawesome.com) to get a single line of code that loads everything automatically.

For detailed installation options, see [Installation Guide](references/installation.md).

## Core Concepts

Web Awesome components are custom HTML elements. They work like native elements but with enhanced functionality.

- **Attributes & Properties**: Configure components via HTML attributes or JavaScript properties
- **Events**: Listen to custom events prefixed with `wa-` (e.g., `wa-change`, `wa-input`)
- **Methods**: Call methods programmatically (e.g., `element.focus()`)
- **Slots**: Insert content into named slots (e.g., `<wa-icon slot="start">`)
- **CSS Parts**: Style internal elements using `::part()` selectors
- **CSS Custom Properties**: Customize appearance with CSS variables

**Important**: Always use closing tags. Custom elements cannot self-close.

```html
<!-- Correct -->
<wa-input></wa-input>

<!-- Incorrect - will not work -->
<wa-input />
```

For complete usage details, see [Usage Guide](references/usage.md).

## Components

### Free Components

#### Actions

- [`<wa-button>`](references/components/button.md) - Buttons represent actions the user can take, such as submitting a form, opening a dialog, or navigating to
another page. ([docs](https://webawesome.com/docs/components/button))
- [`<wa-button-group>`](references/components/button-group.md) - Button groups combine related buttons into a single visual unit. Use them for toolbars, segmented controls,
or any set of actions that belong together. ([docs](https://webawesome.com/docs/components/button-group))
- [`<wa-copy-button>`](references/components/copy-button.md) - Copy buttons copy text to the clipboard when the user activates them. They provide built-in success and
error feedback so users know the copy worked. ([docs](https://webawesome.com/docs/components/copy-button))
- [`<wa-dropdown>`](references/components/dropdown.md) - Dropdowns display a list of options triggered by a button or other element. They support keyboard
navigation, submenus, and checkable items for building menus and context actions. ([docs](https://webawesome.com/docs/components/dropdown))
- [`<wa-dropdown-item>`](references/components/dropdown-item.md) - Dropdown items represent selectable entries within a dropdown menu, including standard actions, checkable
items, and submenu triggers. ([docs](https://webawesome.com/docs/components/dropdown-item))
- [`<wa-qr-code>`](references/components/qr-code.md) - QR codes encode a URL or other short text into a scannable image, rendered client-side using the Canvas API.
Use them to share links, contact info, or Wi-Fi credentials that visitors can scan with a phone. ([docs](https://webawesome.com/docs/components/qr-code))

#### Feedback & Status

- [`<wa-badge>`](references/components/badge.md) - Badges draw attention to adjacent content by displaying a status, count, or label. Use them to highlight
notifications, categorize items, or flag new activity. ([docs](https://webawesome.com/docs/components/badge))
- [`<wa-callout>`](references/components/callout.md) - Callouts display important messages inline with surrounding content. Use them to highlight tips, warnings,
errors, or other information users should not miss. ([docs](https://webawesome.com/docs/components/callout))
- [`<wa-progress-bar>`](references/components/progress-bar.md) - Progress bars show how far along an ongoing operation is as a horizontal fill. Use them for file uploads,
multi-step flows, or any task with measurable progress. ([docs](https://webawesome.com/docs/components/progress-bar))
- [`<wa-progress-ring>`](references/components/progress-ring.md) - Progress rings show how far along a determinate operation is using a circular indicator. Use them as a
compact alternative to progress bars when horizontal space is limited. ([docs](https://webawesome.com/docs/components/progress-ring))
- [`<wa-skeleton>`](references/components/skeleton.md) - Skeletons show placeholder shapes where content will appear once it finishes loading, reducing perceived
wait time and preventing layout shift. ([docs](https://webawesome.com/docs/components/skeleton))
- [`<wa-spinner>`](references/components/spinner.md) - Spinners indicate that an operation is in progress when the duration is unknown. Use them for loading states
where a determinate progress bar isn't practical. ([docs](https://webawesome.com/docs/components/spinner))
- [`<wa-tag>`](references/components/tag.md) - Tags label, categorize, or represent selections with a compact visual marker. Use them for status
indicators, filters, or removable chips. ([docs](https://webawesome.com/docs/components/tag))
- [`<wa-tooltip>`](references/components/tooltip.md) - Tooltips display brief contextual information when the user hovers, focuses, or taps a target element. ([docs](https://webawesome.com/docs/components/tooltip))

#### Form Controls

- [`<wa-checkbox>`](references/components/checkbox.md) - Checkboxes let users toggle an option on or off, or select multiple items from a list. They also support an
indeterminate state for partial selections in groups. ([docs](https://webawesome.com/docs/components/checkbox))
- [`<wa-color-picker>`](references/components/color-picker.md) - Color pickers let users choose a color from a visual palette or by entering a value. They support HEX, RGB,
HSL, and HSV formats with optional alpha channel and swatch presets. ([docs](https://webawesome.com/docs/components/color-picker))
- [`<wa-input>`](references/components/input.md) - Inputs collect single-line data from the user, such as text, numbers, email addresses, and passwords. They
support labels, hints, validation, and prefix or suffix slots. ([docs](https://webawesome.com/docs/components/input))
- [`<wa-number-input>`](references/components/number-input.md) - Number inputs let users enter and edit numeric values, with optional stepper buttons for incrementing and
decrementing. Use them for quantities, measurements, and other numeric form fields. ([docs](https://webawesome.com/docs/components/number-input))
- [`<wa-option>`](references/components/option.md) - Options represent the individual choices inside a select or similar form control. Each option holds a value
and the label shown to the user. ([docs](https://webawesome.com/docs/components/option))
- [`<wa-radio>`](references/components/radio.md) - Radios represent a single option within a mutually exclusive set. Use them inside a radio group when users
must pick exactly one choice from a small list. ([docs](https://webawesome.com/docs/components/radio))
- [`<wa-radio-group>`](references/components/radio-group.md) - Radio groups wrap a set of radios so they function as a single form control with one shared value. They
handle keyboard navigation, labeling, and validation for the group as a whole. ([docs](https://webawesome.com/docs/components/radio-group))
- [`<wa-rating>`](references/components/rating.md) - Ratings display a numeric score as a row of selectable symbols, typically stars. Use them to capture quick
feedback or show an average rating for a product or piece of content. ([docs](https://webawesome.com/docs/components/rating))
- [`<wa-select>`](references/components/select.md) - Selects let users choose one or more values from a dropdown list of predefined options. Use them in forms
when a fixed set of choices needs to fit in limited space. ([docs](https://webawesome.com/docs/components/select))
- [`<wa-slider>`](references/components/slider.md) - Sliders let users choose a numeric value within a defined range by dragging a thumb along a track. ([docs](https://webawesome.com/docs/components/slider))
- [`<wa-switch>`](references/components/switch.md) - Switches toggle a single setting on or off and apply the change immediately, without requiring a form
submission. ([docs](https://webawesome.com/docs/components/switch))
- [`<wa-textarea>`](references/components/textarea.md) - Textareas collect multi-line text input from the user, with optional resizing and character counting. ([docs](https://webawesome.com/docs/components/textarea))

#### Imagery

- [`<wa-animated-image>`](references/components/animated-image.md) - Animated images display GIFs and WEBPs with controls to play and pause them on demand. Use them when you
want motion but need to give users control over when it plays. ([docs](https://webawesome.com/docs/components/animated-image))
- [`<wa-avatar>`](references/components/avatar.md) - Avatars represent a person or object with an image, initials, or icon. Use them in lists, comments, and
profiles to give users visual context at a glance. ([docs](https://webawesome.com/docs/components/avatar))
- [`<wa-carousel>`](references/components/carousel.md) - Carousels display a series of content slides along a horizontal or vertical axis, one or more at a time.
Users can navigate between slides with controls, pagination, or autoplay. ([docs](https://webawesome.com/docs/components/carousel))
- [`<wa-carousel-item>`](references/components/carousel-item.md) - Carousel items represent individual slides within a carousel. ([docs](https://webawesome.com/docs/components/carousel-item))
- [`<wa-comparison>`](references/components/comparison.md) - Comparisons show the visual differences between two pieces of similar content using a draggable divider. Use
them for before/after images, design revisions, or side-by-side previews. ([docs](https://webawesome.com/docs/components/comparison))
- [`<wa-icon>`](references/components/icon.md) - Icons are scalable vector symbols that represent actions, content, or status throughout your application.
They support Font Awesome and custom icon libraries with animation presets. ([docs](https://webawesome.com/docs/components/icon))
- [`<wa-zoomable-frame>`](references/components/zoomable-frame.md) - Zoomable frames embed iframe content with built-in controls for zooming, panning, and managing interaction. ([docs](https://webawesome.com/docs/components/zoomable-frame))

#### Navigation

- [`<wa-breadcrumb>`](references/components/breadcrumb.md) - Breadcrumbs display a trail of links that show users where they are in a site's hierarchy. They help users
understand the current location and navigate back to parent pages. ([docs](https://webawesome.com/docs/components/breadcrumb))
- [`<wa-breadcrumb-item>`](references/components/breadcrumb-item.md) - Breadcrumb items represent individual links inside a breadcrumb, typically one per level of the site
hierarchy. ([docs](https://webawesome.com/docs/components/breadcrumb-item))
- [`<wa-tab>`](references/components/tab.md) - Tabs label and activate an individual panel inside a tab group. ([docs](https://webawesome.com/docs/components/tab))
- [`<wa-tab-group>`](references/components/tab-group.md) - Tab groups organize related content into a single container that displays one panel at a time, with tabs for
switching between them. ([docs](https://webawesome.com/docs/components/tab-group))
- [`<wa-tab-panel>`](references/components/tab-panel.md) - Tab panels hold the content shown for a single tab inside a tab group. ([docs](https://webawesome.com/docs/components/tab-panel))
- [`<wa-tree>`](references/components/tree.md) - Trees allow you to display a hierarchical list of selectable tree items. Items with children can be expanded
and collapsed as desired by the user. ([docs](https://webawesome.com/docs/components/tree))
- [`<wa-tree-item>`](references/components/tree-item.md) - Tree items represent a single hierarchical node inside a tree, and can contain nested items that expand and
collapse. ([docs](https://webawesome.com/docs/components/tree-item))

#### Organization

- [`<wa-card>`](references/components/card.md) - Cards group related content and actions inside a bordered container. Use them to present products, articles,
user profiles, or any self-contained unit of information. ([docs](https://webawesome.com/docs/components/card))
- [`<wa-details>`](references/components/details.md) - Details display a brief summary and expand to reveal additional content. Use them to progressively disclose
information, group related FAQs, or hide advanced options. ([docs](https://webawesome.com/docs/components/details))
- [`<wa-dialog>`](references/components/dialog.md) - Dialogs appear above the page and require the user's immediate attention. Use them for confirmations, forms,
or focused tasks that interrupt the main flow. ([docs](https://webawesome.com/docs/components/dialog))
- [`<wa-divider>`](references/components/divider.md) - Dividers visually separate or group adjacent elements with a horizontal or vertical line. Use them to
establish rhythm and hierarchy within menus, toolbars, and layouts. ([docs](https://webawesome.com/docs/components/divider))
- [`<wa-drawer>`](references/components/drawer.md) - Drawers slide in from the edge of a container to expose additional options and information without
navigating away. Useful for navigation menus, filters, and secondary content. ([docs](https://webawesome.com/docs/components/drawer))
- [`<wa-page>`](references/components/page.md) - Pages scaffold an entire application layout with header, navigation, sidebar, main content, aside, and
footer regions. Use them to structure full pages with minimal markup and responsive behavior built in. ([docs](https://webawesome.com/docs/components/page))
- [`<wa-scroller>`](references/components/scroller.md) - Scrollers wrap overflowing content in an accessible container with visual cues that help users recognize and
navigate scrollable regions. ([docs](https://webawesome.com/docs/components/scroller))
- [`<wa-split-panel>`](references/components/split-panel.md) - Split panels display two adjacent panels separated by a draggable divider, letting users resize each side to
suit their workflow. ([docs](https://webawesome.com/docs/components/split-panel))

#### Utilities

- [`<wa-animation>`](references/components/animation.md) - Animate elements declaratively with nearly 100 baked-in presets, or roll your own with custom keyframes.
Powered by the Web Animations API. ([docs](https://webawesome.com/docs/components/animation))
- [`<wa-format-bytes>`](references/components/format-bytes.md) - Formats a number of bytes as a human-readable string with the appropriate unit, such as kB, MB, or GB.
Supports both byte and bit units with configurable locale. ([docs](https://webawesome.com/docs/components/format-bytes))
- [`<wa-format-date>`](references/components/format-date.md) - Formats a date or time for display using the specified locale and options. Powered by the
Intl.DateTimeFormat API for consistent, localized output. ([docs](https://webawesome.com/docs/components/format-date))
- [`<wa-format-number>`](references/components/format-number.md) - Formats a number for display using the specified locale and options, including currency, percent, and unit
styles. Powered by the Intl.NumberFormat API. ([docs](https://webawesome.com/docs/components/format-number))
- [`<wa-include>`](references/components/include.md) - Fetches an external HTML file and embeds its contents inline on the page. Useful for reusing shared markup
like headers, footers, and partials across multiple pages. ([docs](https://webawesome.com/docs/components/include))
- [`<wa-intersection-observer>`](references/components/intersection-observer.md) - Tracks immediate child elements and fires events as they move in and out of view. Useful for lazy loading,
scroll-triggered animations, and viewport-aware interactions. ([docs](https://webawesome.com/docs/components/intersection-observer))
- [`<wa-markdown>`](references/components/markdown.md) - Markdown elements render markdown content as HTML directly in the browser, making it easy to display
user-generated content or documentation without a server-side build step. ([docs](https://webawesome.com/docs/components/markdown))
- [`<wa-mutation-observer>`](references/components/mutation-observer.md) - Mutation observers watch for changes to an element's DOM tree and emit an event when they occur. Provides a
thin, declarative interface to the browser's MutationObserver API. ([docs](https://webawesome.com/docs/components/mutation-observer))
- [`<wa-popover>`](references/components/popover.md) - Popovers display contextual content and interactive elements in a floating panel anchored to a trigger. Use
them for rich tooltips, menus, or any content that appears on demand without navigating away. ([docs](https://webawesome.com/docs/components/popover))
- [`<wa-popup>`](references/components/popup.md) - Popups declaratively anchor one element to another and keep them positioned together as the page scrolls or
resizes. Primarily a low-level building block for popovers, dropdowns, and tooltips. ([docs](https://webawesome.com/docs/components/popup))
- [`<wa-relative-time>`](references/components/relative-time.md) - Relative times display a date as a localized phrase relative to now, such as "3 hours ago" or "in 2 days".
The phrase updates automatically as time passes and respects the user's locale. ([docs](https://webawesome.com/docs/components/relative-time))
- [`<wa-resize-observer>`](references/components/resize-observer.md) - Resize observers watch their slotted elements for size changes and emit an event when they occur. Provides a
thin, declarative interface to the browser's ResizeObserver API. ([docs](https://webawesome.com/docs/components/resize-observer))


## Themes

Web Awesome includes pre-built themes. Apply a theme by adding its class to the `<html>` element.

### Free Themes
- **Default** - The foundational theme
- **Awesome** - Bright, vibrant color palette
- **Shoelace** - Classic Shoelace styling

### Pro Themes
- **Active** - Green branding with rudimentary palette
- **Brutalist** - Blue branding with default palette
- **Glossy** - Indigo accents with elegant palette
- **Matter** - Purple branding with mild palette
- **Mellow** - Blue branding with natural palette
- **Playful** - Purple branding with rudimentary palette
- **Premium** - Cyan branding with anodized palette
- **Tailspin** - Indigo accents with vogue palette

See [Themes Reference](references/themes.md) for usage details.

## Color Palettes

Each palette provides 10 color hues with 11 tints each.

### Free Palettes
- Default, Bright, Shoelace

### Pro Palettes
- Rudimentary, Elegant, Mild, Natural, Anodized, Vogue

See [Themes Reference](references/themes.md) for palette usage.

## Utilities

Web Awesome provides CSS utilities for common styling tasks:

- **Layout**: `wa-stack`, `wa-cluster`, `wa-grid`, `wa-split`, `wa-flank`, `wa-frame`
- **Spacing**: `wa-gap-*` utilities
- **Text**: Typography utilities
- **Color**: Color variant utilities
- **Rounding**: `wa-border-radius-*` utilities
- **Accessibility**: `wa-visually-hidden` utilities
- **FOUCE Prevention**: `wa-cloak` utility
- **Native Styles**: Enhanced styling for native HTML elements

See [Layout Utilities](references/utilities/layout.md), [Rounding](references/utilities/rounding.md), [Visually Hidden](references/utilities/visually-hidden.md), [FOUCE](references/utilities/fouce.md), and [Native Styles](references/utilities/native.md).

## Design Tokens

Web Awesome uses CSS custom properties (design tokens) for consistent theming:

- **Borders**: `--wa-border-*` for width, radius, style
- **Color**: `--wa-color-*` for surfaces, text, semantic colors
- **Space**: `--wa-space-*` for consistent spacing
- **Typography**: `--wa-font-*` for font families, sizes, weights
- **Shadows**: `--wa-shadow-*` for elevation
- **Focus**: `--wa-focus-*` for focus ring styles
- **Transitions**: `--wa-transition-*` for animation timing

See [Design Tokens](references/tokens/) for full reference.

## Form Controls

Web Awesome form controls are form-associated custom elements supporting native form validation and the Constraint Validation API.

- Use `required`, `pattern`, `minlength`, `maxlength` attributes
- Use `setCustomValidity()` for custom error messages
- Style validation states with `:state(valid)`, `:state(invalid)`, etc.

See [Form Controls Reference](references/form-controls.md) for details.

## Icons

Font Awesome is the default icon library. Use `<wa-icon>` with Font Awesome icon names:

```html
<wa-icon name="house"></wa-icon>
<wa-icon name="gear"></wa-icon>
<wa-icon name="check"></wa-icon>
```

## Framework Integration

Web Awesome works with any framework:

- **React 19+**: Native custom element support with TypeScript types
- **React 18 and below**: Use provided React wrappers
- **Vue**: Works out of the box
- **Angular**: Works out of the box
- **Svelte**: Works out of the box

See framework-specific guides in [references/frameworks/](references/frameworks/).

## Pro Features

[Web Awesome Pro](https://webawesome.com/purchase) includes:

- Pro Components (Data Grid, Date Picker, Rich Text Editor, etc.)
- Pro Themes and Color Palettes
- Theme Builder tool
- Official Figma Design Kit
- Responsive Layout Tools
- Pattern Library
- Priority Support

## Support

- **GitHub Issues**: https://github.com/shoelace-style/webawesome/issues
- **GitHub Discussions**: https://github.com/shoelace-style/webawesome/discussions
- **Discord**: Community chat and support
- **Email**: For account and billing questions

See [Support Reference](references/support.md) for more details.

## Reference Documentation

- [Installation Guide](references/installation.md)
- [Usage Guide](references/usage.md)
- [Form Controls](references/form-controls.md)
- [Customizing](references/customizing.md)
- [Localization](references/localization.md)
- [Themes & Palettes](references/themes.md)
- [Layout Utilities](references/utilities/layout.md)
- [Native Styles](references/utilities/native.md)
- [Design Tokens](references/tokens/) - Borders, Color, Space, Typography, Shadows, Focus, Transitions
- [Framework Guides](references/frameworks/)
