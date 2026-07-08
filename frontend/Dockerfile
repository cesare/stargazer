FROM node:24

USER node
WORKDIR /app
COPY --chown=node:node ./package*.json /app
RUN npm install
