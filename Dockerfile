FROM node:21 AS build

WORKDIR /app

COPY package.json ./
COPY package-lock.json ./
RUN npm install
COPY . ./
RUN npm run build

FROM nginx:1.26-alpine
COPY --from=build /app/build /usr/share/nginx/html