export const checkCompat = (): [boolean, string[]] => {
  const missingFeatures: string[] = [];

  if (window.CustomElementRegistry === undefined) {
    missingFeatures.push("CustomElementRegistry");
  }

  if (window.ElementInternals === undefined) {
    missingFeatures.push("ElementInternals");
  }

  if (window.CustomStateSet === undefined) {
    missingFeatures.push("CustomStateSet");
  }

  return [missingFeatures.length === 0, missingFeatures];
};
