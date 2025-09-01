export const loadNativeStyles = async () => {
  const NativeStyles = (await import("./NativeStyles.svelte")).default;

  new NativeStyles({
    target: document.body,
  });
};
