import type { Attachment } from "svelte/attachments";

export const value = (value: string | number | undefined): Attachment => (
  node: Element) => {
  const update = (value: string | number | undefined) => {
    node.setAttribute("value", value?.toString() ?? "");
  };

  update(value);

  // return {
  //   update,
  // };
};

export const name = (value: string | number | undefined): Attachment => (node: Element) => {
  node.addEventListener("sl-invalid", (e) => e.preventDefault());

  node.setAttribute("name", value?.toString() ?? "");
};

export const checked = (value: boolean | undefined): Attachment => (node: Element) => {
  const update = (value: boolean | undefined) => {
    if (value) {
      node.setAttribute("checked", "");
    } else {
      node.removeAttribute("checked");
    }
  };

  update(value);

  //  return {
  //    update,
  //  };
};
