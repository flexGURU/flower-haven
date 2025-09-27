// mypreset.ts
import { definePreset } from '@primeng/themes';
import Aura from '@primeng/themes/aura';

const MagentaPinkPreset = definePreset(Aura, {
  semantic: {
    primary: {
      50:  '#fff0f7', // very light pink tint
      100: '#ffd9eb',
      200: '#ffb3d7',
      300: '#ff8cc3',
      400: '#f45ba4',
      500: '#e4378d',
      600: '#cc1e77', // <-- main brand color
      700: '#a3185f',
      800: '#7a1247',
      900: '#510b2f',
      950: '#2b0419', // deep accent
    },
  },
});

export default MagentaPinkPreset;
