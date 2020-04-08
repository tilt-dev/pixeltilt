FROM node:13.7.0-alpine3.10

WORKDIR /usr/src/app

ADD package.json package-lock.json ./
RUN npm install

ADD . .

EXPOSE 3000
CMD ["npm", "run", "dev"]
