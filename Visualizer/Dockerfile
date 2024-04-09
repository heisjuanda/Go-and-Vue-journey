FROM node:18.17.1-alpine

RUN npm install -g http-server

RUN mkdir /app

WORKDIR /app

COPY ./package*.json ./

ENV VITE_EMAIL_INDEXER_USER=admin
ENV VITE_EMAIL_INDEXER_PASSWORD=Complexpass#123

RUN npm install

COPY . .

RUN npm run build

CMD [ "http-server", "dist" ]

## Commands
# docker build -t visualizer-truora .

# docker images
# docker run -p 5000:5173 visualizer-truora (localhost:5000)