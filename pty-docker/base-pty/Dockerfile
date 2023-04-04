FROM node:18-alpine
WORKDIR /app
RUN apk add --no-cache --update bash
# <src> 是一个目录，则将目录下的所有文件写入<dest>中
COPY ./build ./build
COPY package*.json .
COPY .env .
RUN apk add --no-cache --update make python3 build-base
RUN npm install --omit=dev

EXPOSE 8081
CMD [ "npm", "run", "start-prod" ]