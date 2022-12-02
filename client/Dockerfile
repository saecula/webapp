FROM node:14-alpine as build
WORKDIR /usr/src/client

COPY package*.json ./

RUN npm install

COPY ./src ./src
COPY ./public ./public
COPY .env* ./

RUN npm run-script build

# production environment
FROM nginx:stable-alpine
COPY --from=build /usr/src/client/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]