// mypreset.ts
import { definePreset } from '@primeng/themes';
import Aura from '@primeng/themes/aura';

const LightPinkPreset = definePreset(Aura, {
  semantic: {
    primary: {
      50: '#fff5f7',  // very soft blush
      100: '#ffe9ed',
      200: '#ffd4da',
      300: '#ffbfc9',
      400: '#ffa9b8',
      500: '#f75d7aff', // base pink
      600: '#ff7a93', // slightly stronger, still soft
      700: '#ff6280',
      800: '#ff4a6e',
      900: '#ff345d', // most saturated but not dark
      950: '#ff1b4c', // deep accent but still bright
    },
  },
});

export default LightPinkPreset;
