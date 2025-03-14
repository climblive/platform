import js from "@eslint/js";
import eslintPluginSvelte from "eslint-plugin-svelte";
import globals from "globals";
import svelteParser from "svelte-eslint-parser";
import tsEslint from "typescript-eslint";

export default [
  js.configs.recommended,
  ...tsEslint.configs.strict,
  ...eslintPluginSvelte.configs["flat/recommended"],
  {
    rules: {
      "no-console": "error",
    },
  },
  {
    files: ["**/*.ts"],
    languageOptions: {
      parser: tsEslint.parser,
    },
  },
  {
    files: ["**/*.svelte"],
    languageOptions: {
      parser: svelteParser,
      parserOptions: {
        parser: tsEslint.parser,
      },
      globals: {
        ...globals.browser,
        $$Generic: 'readonly',
        NodeListOf: 'readonly',
      },
    },
    rules: {
      "svelte/valid-compile": "off",
    },
  },
];
