// Runtime/build-time configuration derived from environment.
//
// NEXT_PUBLIC_* vars are baked into the client bundle at build time, so
// changing them requires a dashboard rebuild (see docker-compose args).

/**
 * Port the proxy server listens on, shown in copy-able proxy URLs.
 * Mirrors the backend PROXY_PORT env var. Defaults to 8000 to match the
 * stack defaults. Previously hard-coded in the Proxy Users page (issue #30).
 */
export const PROXY_PORT = process.env.NEXT_PUBLIC_PROXY_PORT || "8000"
