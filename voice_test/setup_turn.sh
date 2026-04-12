#!/bin/bash

echo "🚀 AnCord TURN Server Setup"
echo "==========================="
echo ""

# Проверяем наличие curl
if ! command -v curl &> /dev/null; then
    echo "❌ curl не установлен. Устанавливаем..."
    sudo apt update && sudo apt install -y curl
    echo "✅ curl установлен"
fi

# Получаем публичный IP
echo "🌐 Получаем публичный IP..."
SERVICES=(
    "https://ifconfig.me"
    "https://api.ipify.org"
    "https://checkip.amazonaws.com"
)

for service in "${SERVICES[@]}"; do
    echo "Пробуем $service..."
    IP=$(curl -s --max-time 5 "$service" 2>/dev/null)
    if [ $? -eq 0 ] && [ -n "$IP" ]; then
        echo ""
        echo "✅ Ваш публичный IP: $IP"
        break
    fi
done

if [ -z "$IP" ]; then
    echo "❌ Не удалось получить публичный IP"
    echo "Используем 127.0.0.1 для локального тестирования"
    IP="127.0.0.1"
fi

# Создаем .env файл если не существует
if [ ! -f "../.env" ]; then
    echo ""
    echo "📝 Создаем .env файл..."
    cat > ../.env << EOF
PORT=8080

# TURN сервер настройки (обязательно)
TURN_PUBLIC_IP=$IP
TURN_PORT=3478
TURN_USERNAME=ancord-user
TURN_PASSWORD=secure-password-$(date +%s)
EOF
    echo "✅ .env файл создан"
else
    echo ""
    echo "📝 Обновляем TURN_PUBLIC_IP в .env файле..."
    # Проверяем, есть ли уже TURN настройки
    if grep -q "TURN_PUBLIC_IP" ../.env; then
        sed -i "s/TURN_PUBLIC_IP=.*/TURN_PUBLIC_IP=$IP/" ../.env
    else
        # Добавляем TURN настройки в конец файла
        echo "" >> ../.env
        echo "# TURN сервер настройки (обязательно)" >> ../.env
        echo "TURN_PUBLIC_IP=$IP" >> ../.env
        echo "TURN_PORT=3478" >> ../.env
        echo "TURN_USERNAME=ancord-user" >> ../.env
        echo "TURN_PASSWORD=secure-password-$(date +%s)" >> ../.env
    fi
    echo "✅ .env файл обновлен"
fi

echo ""
echo "🎉 Настройка завершена!"
echo ""
echo "📋 Конфигурация:"
echo "   TURN_PUBLIC_IP: $IP"
echo "   TURN_PORT: 3478"
echo ""
echo "🚀 Запуск сервера:"
echo "   cd voice_test && go run ./sfu"
echo ""
echo "🐳 Или с Docker:"
echo "   docker-compose up -d"
echo ""
echo "💡 Для продакшена замените 127.0.0.1 на реальный публичный IP"
echo ""

# Проверяем наличие curl
if ! command -v curl &> /dev/null; then
    echo "❌ curl не установлен. Устанавливаем..."
    sudo apt update && sudo apt install -y curl
    echo "✅ curl установлен"
fi

# Получаем публичный IP
echo "🌐 Получаем публичный IP..."
SERVICES=(
    "https://ifconfig.me"
    "https://api.ipify.org"
    "https://checkip.amazonaws.com"
)

for service in "${SERVICES[@]}"; do
    echo "Пробуем $service..."
    IP=$(curl -s --max-time 5 "$service" 2>/dev/null)
    if [ $? -eq 0 ] && [ -n "$IP" ]; then
        echo ""
        echo "✅ Ваш публичный IP: $IP"
        break
    fi
done

if [ -z "$IP" ]; then
    echo "❌ Не удалось получить публичный IP"
    echo "Укажите его вручную в .env файле"
    exit 1
fi

# Создаем .env файл если не существует
if [ ! -f "../.env" ]; then
    echo ""
    echo "📝 Создаем .env файл..."
    cat > ../.env << EOF
PORT=8080

# TURN сервер настройки (обязательно)
TURN_PUBLIC_IP=$IP
TURN_PORT=3478
TURN_USERNAME=ancord-user
TURN_PASSWORD=secure-password-$(date +%s)
EOF
    echo "✅ .env файл создан"
else
    echo ""
    echo "📝 Обновляем TURN_PUBLIC_IP в .env файле..."
    # Проверяем, есть ли уже TURN настройки
    if grep -q "TURN_PUBLIC_IP" ../.env; then
        sed -i "s/TURN_PUBLIC_IP=.*/TURN_PUBLIC_IP=$IP/" ../.env
    else
        # Добавляем TURN настройки в конец файла
        echo "" >> ../.env
        echo "# TURN сервер настройки (обязательно)" >> ../.env
        echo "TURN_PUBLIC_IP=$IP" >> ../.env
        echo "TURN_PORT=3478" >> ../.env
        echo "TURN_USERNAME=ancord-user" >> ../.env
        echo "TURN_PASSWORD=secure-password-$(date +%s)" >> ../.env
    fi
    echo "✅ .env файл обновлен"
fi

echo ""
echo "🎉 Настройка завершена!"
echo ""
echo "📋 Конфигурация:"
echo "   TURN_PUBLIC_IP: $IP"
echo "   TURN_PORT: 3478"
echo ""
echo "🚀 Запуск сервера:"
echo "   go run ./sfu"
echo ""
echo "🐳 Или с Docker:"
echo "   docker-compose up -d"