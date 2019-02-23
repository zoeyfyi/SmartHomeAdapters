FROM node:8 as react-build
WORKDIR /app
COPY . ./
RUN npm install
RUN npm run build

FROM nginx:alpine
COPY --from=react-build /app/build /usr/share/nginx/html
EXPOSE 80
CMD nginx -g 'daemon off;'
