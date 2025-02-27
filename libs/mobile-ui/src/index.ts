import { camelCase } from "change-case";
import convert, { HSL } from "color-convert";

import { colors as colorsHSL } from "./theme/colors";

export * from "./components";

const hslStrToRgbStr = (hslStr: string) =>
  `#${convert.hsl.hex(
    ...(hslStr.replace(/%/g, "").split(" ").map(Number) as HSL),
  )}`;

const colors = Object.entries(colorsHSL).reduce(
  (acc, [key, value]) => ({
    ...acc,
    [camelCase(key)]: Object.entries(value).reduce(
      (a, [k, v]) => ({
        ...a,
        [camelCase(k)]: hslStrToRgbStr(v),
      }),
      {},
    ),
  }),
  {},
);

console.log(colors);
