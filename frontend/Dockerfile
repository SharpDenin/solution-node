FROM node:20

WORKDIR /app

COPY package*.json ./
RUN npm install
RUN npm install axios jwt-decode react-router-dom

COPY . .

EXPOSE 3000

CMD ["npm", "run", "dev", "--", "--host"]