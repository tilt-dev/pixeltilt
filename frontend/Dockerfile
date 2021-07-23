FROM node:16-alpine3.14

WORKDIR /usr/src/app

ADD package.json package-lock.json ./
RUN npm install

ADD . .

EXPOSE 3000
CMD ["npm", "run", "dev"]
