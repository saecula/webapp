FROM node:14-alpine
WORKDIR /usr/src/gogo

COPY package*.json ./

RUN npm install

COPY ./src ./src
COPY ./public ./public

RUN npm run-script build

EXPOSE 3000

CMD ["npm", "start"]