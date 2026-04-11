#!/bin/bash

echo "🚀 AnCord - Полный запуск проекта"
echo "=================================="
echo ""

# Проверяем, что мы в корне проекта
if [ ! -f ".env" ]; then
    echo "❌ .env файл не найден в корне проекта"
    echo "Запустите: cd voice_test && ./setup_turn.sh"
    exit 1
fi

echo "✅ Конфигурация найдена"

# Запускаем voice_test сервер в фоне
echo "🎯 Запуск SFU сервера..."
cd voice_test
go run ./sfu &
SFU_PID=$!

# Ждем немного
sleep 2

# Проверяем, что сервер запустился
if kill -0 $SFU_PID 2>/dev/null; then
    echo "✅ SFU сервер запущен (PID: $SFU_PID)"
else
    echo "❌ Ошибка запуска SFU сервера"
    exit 1
fi

# Запускаем HTTP сервер для тестов
echo "🌐 Запуск HTTP сервера для тестов..."
cd ..
python3 -m http.server 8000 &
HTTP_PID=$!

echo ""
echo "🎉 Серверы запущены!"
echo ""
echo "📋 Доступ:"
echo "   SFU/WebSocket: http://localhost:8080"
echo "   TURN: UDP 3478"
echo "   Тестовая страница: http://localhost:8000"
echo ""
echo "🛑 Для остановки нажмите Ctrl+C"

# Функция очистки при выходе
cleanup() {
    echo ""
    echo "🧹 Останавливаем серверы..."
    kill $SFU_PID 2>/dev/null
    kill $HTTP_PID 2>/dev/null
    echo "✅ Серверы остановлены"
    exit 0
}

# Ожидаем сигнала прерывания
trap cleanup SIGINT SIGTERM
wait