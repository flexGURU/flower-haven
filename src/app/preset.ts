// mypreset.ts
import { definePreset } from '@primeng/themes';
import Aura from '@primeng/themes/aura';

const MagentaPinkPreset = definePreset(Aura, {
  semantic: {
    primary: {
      50:  '#fdf2f8', // barely there pink - perfect for hover backgrounds
      100: '#fce7f3', // subtle highlight backgrounds
      200: '#fbcfe8', // light borders and dividers
      300: '#f9a8d4', // disabled states, secondary elements
      400: '#f472b6', // hover states on buttons
      500: '#ec4899', // <-- MAIN brand color - vibrant, professional pink
      600: '#db2777', // active/pressed states - deeper, richer
      700: '#be185d', // dark mode primary, strong accents
      800: '#9d174d', // text on light backgrounds
      900: '#831843', // headings, emphasis
      950: '#500724', // deep contrast for borders/shadows
    },
  },
});

export default MagentaPinkPreset;