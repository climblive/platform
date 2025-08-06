import type { Attachment } from "svelte/attachments";

export const value =
  (value: string | number | undefined): Attachment =>
    (node: Element) => {
      node.setAttribute("value", value?.toString() ?? "");
    };

export const name =
  (value: string | number | undefined): Attachment =>
    (node: Element) => {
      node.addEventListener("wa-invalid", (e) => e.preventDefault());

      node.setAttribute("name", value?.toString() ?? "");
    };
