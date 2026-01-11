type Feature =
  | "CustomElementRegistry"
  | "ElementInternals"
  | "CustomStateSet"
  | "@supports selector(&)";

export const checkCompat = (): [boolean, Feature[]] => {
  const missingFeatures: Feature[] = [];

  if (window.CustomElementRegistry === undefined) {
    missingFeatures.push("CustomElementRegistry");
  }

  if (window.ElementInternals === undefined) {
    missingFeatures.push("ElementInternals");
  }

  if (window.CustomStateSet === undefined) {
    missingFeatures.push("CustomStateSet");
  }

  if (CSS.supports("selector(&)") === false) {
    missingFeatures.push("@supports selector(&)");
  }

  return [missingFeatures.length === 0, missingFeatures];
};
