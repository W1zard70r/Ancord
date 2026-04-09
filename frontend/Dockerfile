FROM node:22-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Сервер (Nginx) для статики
FROM nginx:alpine
# Копируем собранные файлы из первого этапа в nginx
COPY --from=builder /app/dist /usr/share/nginx/html
# Простой конфиг для SPA (чтобы роутинг React работал)
RUN echo "server { listen 80; location / { root /usr/share/nginx/html; index index.html; try_files \$uri \$uri/ /index.html; } }" > /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]