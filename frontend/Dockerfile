FROM node:23-alpine AS base

ARG PUBLIC_API_URL
ENV PUBLIC_API_URL=$PUBLIC_API_URL

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /app
COPY . /app

FROM base AS prod-deps
RUN pnpm install --prod --frozen-lockfile

FROM base AS build
RUN pnpm install --frozen-lockfile && pnpm run build

FROM node:23-alpine AS prod

ENV PORT=8080
EXPOSE 8080

RUN apk add curl

USER 10001:10001
WORKDIR /app
COPY --from=prod-deps /app/node_modules /app/node_modules
COPY --from=build /app/build /app/build
COPY package.json /app

CMD [ "node", "./build/index.js" ]
