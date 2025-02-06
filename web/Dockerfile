FROM node:20 AS base
RUN corepack --version
RUN corepack enable
COPY . /app
WORKDIR /app

RUN pnpm i
RUN pnpm --filter=* build

FROM nginx
COPY --from=base /app/scorecard/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf