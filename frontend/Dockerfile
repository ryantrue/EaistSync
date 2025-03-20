# Используем официальный образ nginx на базе Alpine Linux
FROM nginx:alpine

# Устанавливаем рабочую директорию для упрощения работы с путями
WORKDIR /usr/share/nginx/html

# Копируем папку static в рабочую директорию (сохраняя структуру)
COPY static ./static

# Удаляем стандартный index.html и создаём символьную ссылку на ваш index.html
RUN rm -f index.html && \
    ln -s static/index.html index.html

# Копируем кастомную конфигурацию Nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Открываем порт 80 для входящих соединений
EXPOSE 80

# (Опционально) Добавляем healthcheck для мониторинга доступности сервера
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost/ || exit 1