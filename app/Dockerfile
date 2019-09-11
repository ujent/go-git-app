FROM node:8.16.1 as build

WORKDIR /usr/src/git-app

COPY package.json package-lock.json ./

RUN npm install

COPY . .

RUN npm run build

FROM nginx:1.17.3

COPY --from=build /usr/src/git-app/build /usr/share/nginx/html