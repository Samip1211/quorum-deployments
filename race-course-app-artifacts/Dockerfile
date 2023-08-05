FROM node:14-bullseye

WORKDIR home

RUN mkdir server && mkdir client && mkdir contracts


COPY ./server /home/server
COPY ./client /home/client
COPY ./contracts /home/contracts
COPY package.json /home

ENTRYPOINT [ "npm", "run", "start" ]