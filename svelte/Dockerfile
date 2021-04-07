FROM node:14-alpine
WORKDIR /usr/src/svelte

COPY rollup.config.js ./
COPY package*.json ./

RUN npm install

COPY ./src ./src
COPY ./public ./public
COPY .env .

RUN npm run-script build

EXPOSE 5000

ENV HOST=0.0.0.0

CMD ["npm", "start"]