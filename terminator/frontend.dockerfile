FROM node:13.7.0-alpine3.10

WORKDIR /usr/src/app

ADD frontend/package.json frontend/package-lock.json ./
RUN npm install

ADD frontend .

EXPOSE 3000
CMD ["npm", "run", "dev"]
