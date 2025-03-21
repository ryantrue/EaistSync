# Этап 1: Сборка приложения
FROM node:23-alpine as build

WORKDIR /app

# Копируем файлы зависимостей и устанавливаем их
COPY package*.json ./
RUN npm install

# Копируем исходный код приложения
COPY . .

# Собираем приложение в production-режиме (результат окажется в папке build)
RUN npm run build

# Этап 2: Раздача статики через nginx
FROM nginx:alpine

# Задаём рабочую директорию
WORKDIR /usr/share/nginx/html

# Удаляем дефолтное содержимое директории
RUN rm -rf ./*

# Копируем собранное React-приложение из предыдущего этапа
COPY --from=build /app/build .

# Если у вас есть кастомный конфиг nginx, копируйте его
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Открываем порт 80
EXPOSE 80

# Healthcheck для мониторинга работоспособности контейнера
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost/ || exit 1

CMD ["nginx", "-g", "daemon off;"]