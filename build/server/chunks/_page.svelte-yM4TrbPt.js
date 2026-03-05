import { b as ssr_context } from './context-CXhJZien.js';

function onDestroy(fn) {
  /** @type {SSRContext} */
  ssr_context.r.on_destroy(fn);
}
function stopPolling() {
}
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    onDestroy(() => stopPolling());
    $$renderer2.push(`<div class="min-h-screen bg-zinc-950 text-zinc-100 font-sans antialiased svelte-1uha8ag"><header class="sticky top-0 z-10 bg-zinc-950/80 backdrop-blur border-b border-zinc-800 px-6 py-3 flex items-center justify-between svelte-1uha8ag"><div class="flex items-center gap-3 svelte-1uha8ag"><svg class="w-6 h-6 text-sky-400 svelte-1uha8ag" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M9 3H5a2 2 0 00-2 2v4m6-6h10a2 2 0 012 2v4M9 3v18m0 0h10a2 2 0 002-2V9M9 21H5a2 2 0 01-2-2V9m0 0h18" class="svelte-1uha8ag"></path></svg> <span class="text-lg font-semibold tracking-tight svelte-1uha8ag">Host Monitor</span> `);
    {
      $$renderer2.push("<!--[!-->");
    }
    $$renderer2.push(`<!--]--></div> <div class="flex items-center gap-2 svelte-1uha8ag">`);
    {
      $$renderer2.push("<!--[-->");
      $$renderer2.push(`<span class="status-badge status-loading svelte-1uha8ag">Connecting…</span>`);
    }
    $$renderer2.push(`<!--]--> `);
    {
      $$renderer2.push("<!--[!-->");
    }
    $$renderer2.push(`<!--]--></div></header> <main class="px-4 sm:px-6 py-6 max-w-6xl mx-auto space-y-6 svelte-1uha8ag">`);
    {
      $$renderer2.push("<!--[-->");
      $$renderer2.push(`<div class="card-base flex items-center justify-center gap-3 py-16 text-zinc-400 svelte-1uha8ag"><svg class="w-5 h-5 animate-spin text-sky-400 svelte-1uha8ag" fill="none" viewBox="0 0 24 24"><circle class="opacity-25 svelte-1uha8ag" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75 svelte-1uha8ag" fill="currentColor" d="M4 12a8 8 0 018-8v8z"></path></svg> <span class="svelte-1uha8ag">Loading host information…</span></div>`);
    }
    $$renderer2.push(`<!--]--></main> <footer class="text-center text-zinc-600 text-xs py-6 px-4 svelte-1uha8ag">`);
    {
      $$renderer2.push("<!--[!-->");
    }
    $$renderer2.push(`<!--]--></footer></div>`);
  });
}

export { _page as default };
//# sourceMappingURL=_page.svelte-yM4TrbPt.js.map
