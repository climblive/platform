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

export const checked =
  (value: boolean | undefined): Attachment =>
  (node: Element) => {
    if (value) {
      node.setAttribute("checked", "");
    } else {
      node.removeAttribute("checked");
    }
  };
