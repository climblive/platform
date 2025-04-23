export const value = (
  node: HTMLElement,
  value: string | number | undefined,
) => {
  const update = (value: string | number | undefined) => {
    node.setAttribute("value", value?.toString() ?? "");
  }

  update(value);

  return {
    update
  }
};

export const name = (node: HTMLElement, value: string | number | undefined) => {
  node.addEventListener("sl-invalid", (e) => e.preventDefault());

  node.setAttribute("name", value?.toString() ?? "");
};

export const checked = (node: HTMLElement, value: boolean | undefined) => {
  const update = (value: boolean | undefined) => {
    if (value) {
      node.setAttribute("checked", "");
    } else {
      node.removeAttribute("checked");
    }
  }

  update(value);

  return {
    update
  }
};
