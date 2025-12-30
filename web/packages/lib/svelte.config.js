import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

export default {
  // Consult https://svelte.dev/docs#compile-time-svelte-preprocess
  // for more information about preprocessors
  preprocess: vitePreprocess(),
  compilerOptions: {
    warningFilter: (warning) => !ignoreWarning[warning.code]?.(warning),
  },
};

const ignoreWarning = {
  a11y_no_static_element_interactions: (w) => w.message.startsWith("`<wa-"),
  a11y_click_events_have_key_events: () => true,
};
