export const value = (
  node: HTMLElement,
  _: string | number | undefined,
) => {
  return {
    update: (value: string | number | undefined) => {
      node.setAttribute("value", value?.toString() ?? "");
    }
  }
};

export const name = (node: HTMLElement, value: string | number | undefined) => {
  node.addEventListener("sl-invalid", (e) => e.preventDefault());

  $effect(() => {
    node.setAttribute("name", value?.toString() ?? "");
  });
};

export const checked = (node: HTMLElement, value: boolean | undefined) => {
  $effect(() => {
    if (value) {
      node.setAttribute("checked", "");
    } else {
      node.removeAttribute("checked");
    }
  });
};
