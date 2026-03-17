# ADR-001: Next.js over Plain React (Create React App / Vite)

## Status
Accepted

## Date
2026-02-18

## Context
The Coffee Journal needs a frontend framework. The two main options are:

1. **Plain React** (via Vite or CRA) — client-side SPA
2. **Next.js 16** (App Router) — full-stack React framework with SSR/SSG

The app has specific requirements that influence this decision:
- Bean detail pages must be shareable with rich social previews (Open Graph)
- The app is designed to go public in the future, so SEO matters
- Some pages are read-heavy (timeline, bean detail) while others are interactive (forms, voice recorder)
- The backend is a separate Go API, so we don't need Next.js API routes for data logic

## Decision
Use **Next.js 16 (latest stable) with the App Router**.

## Rationale

### Server-Side Rendering for SEO and Social Previews
Bean detail pages need proper `<meta>` tags rendered in the initial HTML for link previews on Twitter, iMessage, Slack, etc. A plain React SPA serves an empty `<div id="root">` — crawlers and link unfurlers won't see any content. Next.js server components render the full HTML on the server, including dynamic Open Graph metadata via `generateMetadata()`.

### Server Components Reduce Client Bundle
Read-heavy pages (timeline, beans list, bean detail, search results) can be server components that ship zero JavaScript to the browser. Only interactive pieces (forms, voice recorder, search bar) need `"use client"`. With a plain React SPA, the entire app ships to the client regardless.

### Built-in Routing
Next.js file-based routing with the App Router eliminates the need for React Router. Nested layouts (`app/beans/[id]/layout.tsx`) map naturally to the app's URL structure. Less dependencies, less configuration.

### Simplified Data Fetching
Server components can `fetch()` directly from the Go API during rendering — no useEffect, no loading spinners for initial page loads, no client-side state management for server data. This removes the need for TanStack Query or SWR for read pages.

### API Proxy via Rewrites
`next.config.ts` rewrites let client components call `/api/*` which proxies to the Go backend. This avoids CORS configuration entirely. A plain React SPA would require either CORS headers on the Go API or a separate proxy setup.

## Alternatives Considered

### Vite + React SPA
- **Pros**: Simpler setup, faster dev server startup, smaller conceptual surface area
- **Cons**: No SSR — social previews won't work without a separate prerender service. Would need React Router, a data fetching library, and CORS configuration. Every page ships JavaScript to the client.

### Remix
- **Pros**: Similar SSR capabilities, good data loading patterns
- **Cons**: Smaller ecosystem, less deployment flexibility. Vercel is purpose-built for Next.js with zero-config deployment. The team has more Next.js experience.

### Astro
- **Pros**: Excellent for content-heavy static sites, ships minimal JS
- **Cons**: Less suited for the interactive parts of the app (forms, voice recording). Would need to integrate a UI framework anyway for client islands.

## Consequences

### Positive
- Social sharing works out of the box with `generateMetadata()`
- Search engines can index all bean pages
- Smaller client bundles for read-only pages
- No need for React Router, TanStack Query, or CORS middleware
- Zero-config deployment to Vercel

### Negative
- More complex mental model (server vs. client components, when to use `"use client"`)
- Slightly slower dev server cold start compared to Vite
- Tighter coupling to Vercel's ecosystem for optimal deployment
- Need to be deliberate about the server/client boundary in each component
