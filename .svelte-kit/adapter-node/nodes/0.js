import * as universal from '../entries/pages/_layout.ts.js';

export const index = 0;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_layout.svelte.js')).default;
export { universal };
export const universal_id = "src/routes/+layout.ts";
export const imports = ["_app/immutable/nodes/0.C6CpbhCs.js","_app/immutable/chunks/BBT1V1ic.js","_app/immutable/chunks/CtuWOtrZ.js","_app/immutable/chunks/Ctl8bicv.js"];
export const stylesheets = ["_app/immutable/assets/0.CvjncEa6.css"];
export const fonts = [];
