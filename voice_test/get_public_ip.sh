#!/bin/bash

echo "=== Получение публичного IP для TURN сервера ==="
echo ""

# Пробуем разные сервисы
SERVICES=(
    "https://ifconfig.me"
    "https://api.ipify.org"
    "https://checkip.amazonaws.com"
    "https://ipinfo.io/ip"
)

for service in "${SERVICES[@]}"; do
    echo "Пробуем $service..."
    IP=$(curl -s --max-time 5 "$service" 2>/dev/null)
    if [ $? -eq 0 ] && [ -n "$IP" ]; then
        echo ""
        echo "✅ Ваш публичный IP: $IP"
        echo ""
        echo "Добавьте в .env файл:"
        echo "TURN_PUBLIC_IP=$IP"
        echo ""
        echo "И перезапустите сервер:"
        echo "go run ./sfu"
        exit 0
    fi
done

echo "❌ Не удалось получить публичный IP"
echo "Проверьте подключение к интернету или укажите IP вручную"
exit 1