FROM node:10-stretch-slim

RUN apt-get -y --allow-unauthenticated update
RUN apt-get -y --allow-unauthenticated install git


WORKDIR /simplejson
RUN git clone https://github.com/bergquist/fake-simple-json-datasource.git
WORKDIR /simplejson/fake-simple-json-datasource
RUN ls -lh
RUN npm install
CMD node index.js
