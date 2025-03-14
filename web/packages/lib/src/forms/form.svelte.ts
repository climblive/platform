export const value = (
  node: HTMLElement,
  value: string | number | undefined,
) => {
  $effect(() => {
    node.setAttribute("value", value?.toString() ?? "");
  });
};

export const name = (
  node: HTMLElement,
  value: string | number | undefined,
) => {
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