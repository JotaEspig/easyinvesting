FROM node:20 AS build

WORKDIR /app
COPY . .

RUN npm install
RUN npm run build -- --configuration production

FROM nginx:stable-alpine

COPY --from=build /app/dist/frontend/browser/* /usr/share/nginx/html/

COPY nginx/default.conf /etc/nginx/conf.d/default.conf

COPY nginx/ssl/server.crt /etc/ssl/certs/server.crt
COPY nginx/ssl/server.key /etc/ssl/private/server.key

EXPOSE 443

CMD ["nginx", "-g", "daemon off;"]
