# Construye el frontend Vue dentro de la imagen del proxy público.
FROM node:22-alpine AS frontend-build
WORKDIR /build/frontend
COPY frontend/package.json frontend/yarn.lock frontend/.yarnrc.yml ./
COPY frontend/.yarn ./.yarn
RUN corepack enable && yarn install --immutable
COPY frontend ./
RUN yarn build-only

FROM nginx:1.27-alpine
COPY deploy/nginx/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=frontend-build /build/internal/web/dist /usr/share/nginx/html
EXPOSE 80
