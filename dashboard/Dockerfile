# Stage 1: Dependencies
FROM node:20-alpine AS deps
WORKDIR /app

# Install pnpm
RUN corepack enable && corepack prepare pnpm@10.19.0 --activate

# Install dependencies based on the preferred package manager.
# --package-import-method=copy avoids pnpm's default hardlinking, which fails on
# filesystems that don't support cross-file hardlinks in the build layer (e.g.
# ZFS) — see issue #23. node_modules is discarded in the multi-stage build, so
# the extra copy cost is negligible.
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile --package-import-method=copy

# Stage 2: Builder
FROM node:20-alpine AS builder
WORKDIR /app

# Install pnpm
RUN corepack enable && corepack prepare pnpm@10.19.0 --activate

COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Set environment variables for build
ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production

# Accept build-time API URL (baked into Next.js static bundle)
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}

# Proxy port shown in copy-able proxy URLs (baked into the bundle at build time)
ARG NEXT_PUBLIC_PROXY_PORT
ENV NEXT_PUBLIC_PROXY_PORT=${NEXT_PUBLIC_PROXY_PORT}

# Build the application
RUN pnpm run build

# Stage 3: Runner
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# Create a non-root user
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy necessary files from builder
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

# Set correct permissions
RUN chown -R nextjs:nodejs /app

USER nextjs

EXPOSE 3000

ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]
