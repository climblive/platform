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
