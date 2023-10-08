FROM node:18-alpine
WORKDIR /app
RUN apk add --no-cache --update bash make python3 build-base
# <src> 是一个目录，则将目录下的所有文件写入<dest>中
COPY ./build ./build
COPY package*.json .
COPY .env .
# make / python3 / build-base是node-pty所需的依赖项
RUN npm install --omit=dev

EXPOSE 3220
CMD [ "npm", "run", "start-prod" ]