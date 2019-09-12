FROM node:8.16.1 as build

WORKDIR /usr/src/git-app

COPY package.json package-lock.json ./

RUN npm install

COPY . .

RUN REACT_APP_BUILD_MODE=release npm run build

FROM nginx:1.17.3

COPY --from=build /usr/src/git-app/build /usr/share/nginx/html

COPY nginx.template /etc/nginx/conf.d/nginx.template

CMD /bin/bash -c "envsubst < /etc/nginx/conf.d/nginx.template > /etc/nginx/conf.d/default.conf && exec nginx -g 'daemon off;'"