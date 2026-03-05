const manifest = (() => {
function __memo(fn) {
	let value;
	return () => value ??= (value = fn());
}

return {
	appDir: "_app",
	appPath: "_app",
	assets: new Set(["robots.txt"]),
	mimeTypes: {".txt":"text/plain"},
	_: {
		client: {start:"_app/immutable/entry/start.DK52nPUw.js",app:"_app/immutable/entry/app.gxzy13n0.js",imports:["_app/immutable/entry/start.DK52nPUw.js","_app/immutable/chunks/DnsSx_Tq.js","_app/immutable/chunks/CtuWOtrZ.js","_app/immutable/chunks/B9a050j0.js","_app/immutable/entry/app.gxzy13n0.js","_app/immutable/chunks/CtuWOtrZ.js","_app/immutable/chunks/DDxd7z1_.js","_app/immutable/chunks/BBT1V1ic.js","_app/immutable/chunks/B9a050j0.js","_app/immutable/chunks/CJcdje4P.js","_app/immutable/chunks/Ctl8bicv.js"],stylesheets:[],fonts:[],uses_env_dynamic_public:false},
		nodes: [
			__memo(() => import('./chunks/0-C088_aTT.js')),
			__memo(() => import('./chunks/1-B2ZsChhU.js')),
			__memo(() => import('./chunks/2-CkN0brWe.js'))
		],
		remotes: {
			
		},
		routes: [
			{
				id: "/",
				pattern: /^\/$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 2 },
				endpoint: null
			},
			{
				id: "/api/hostinfo",
				pattern: /^\/api\/hostinfo\/?$/,
				params: [],
				page: null,
				endpoint: __memo(() => import('./chunks/_server.ts-C7gAOHfH.js'))
			}
		],
		prerendered_routes: new Set([]),
		matchers: async () => {
			
			return {  };
		},
		server_assets: {}
	}
}
})();

const prerendered = new Set([]);

const base = "";

export { base, manifest, prerendered };
//# sourceMappingURL=manifest.js.map
