FROM node:23-slim AS base

ARG PUBLIC_API_URL
ENV PUBLIC_API_URL=$PUBLIC_API_URL

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
COPY . /app
WORKDIR /app

FROM base AS prod-deps
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile

FROM base AS build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build

FROM node:23-alpine AS prod
COPY --from=prod-deps /app/node_modules /app/node_modules
COPY --from=build /app/.svelte-kit/output/server /app/.svelte-kit/output/server

ENV PORT=8080
EXPOSE 8080

RUN apk add curl

USER 10001:10001
CMD [ "pnpm", "start" ]