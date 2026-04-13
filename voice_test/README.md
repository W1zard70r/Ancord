# AnCord Voice SFU Server

WebRTC SFU (Selective Forwarding Unit) сервер для голосового общения с собственным TURN сервером через pion/turn.

## Быстрый старт

### 1. Настройка TURN сервера

```bash
# Автоматическая настройка (получение IP + создание .env)
cd voice_test
./setup_turn.sh
```

### 2. Запуск сервера

```bash
# Из директории voice_test
cd voice_test
go run ./sfu

# Или из корня проекта
go run ./voice_test/sfu
```

### 3. Тестирование

```bash
# В другом терминале запустите HTTP сервер
python3 -m http.server 8000
# Откройте http://localhost:8000 в браузере
```

## ⚙️ Настройка

### Переменные окружения (.env в корне проекта)

```env
PORT=8080

# TURN сервер настройки (ОБЯЗАТЕЛЬНО)
TURN_PUBLIC_IP=127.0.0.1  # Для тестирования
TURN_PORT=3478
TURN_USERNAME=ancord-user
TURN_PASSWORD=secure-password
```

### Для продакшена

```bash
# Получите реальный публичный IP
curl ifconfig.me

# Обновите .env
TURN_PUBLIC_IP=ВАШ_ПУБЛИЧНЫЙ_IP
```

## 🐳 Docker развертывание

### Docker:
```bash
cd voice_test
docker build -t ancord-sfu:latest .
docker run -p 8080:8080 -p 3478:3478/udp --env-file ../.env ancord-sfu:latest
```

### Docker Compose:
```bash
docker-compose up -d
```

## 📊 Архитектура

- **SFU (Selective Forwarding Unit)**: Маршрутизация медиа-потоков
- **TURN сервер (pion/turn)**: Обход NAT и firewall
- **WebSocket сигнализация**: Обмен SDP offer/answer
- **Docker**: Контейнеризация

## 🔒 Безопасность

- TURN credentials генерируются динамически (24 часа)
- WebSocket с проверкой origin
- Переменные окружения для чувствительных данных
- Собственный TURN сервер (pion/turn)

## 📋 Логи сервера

```
TURN сервер запущен на 127.0.0.1:3478
Server started on port 8080
```

## 🛠️ Устранение неполадок

### Проблема: "TURN_PUBLIC_IP обязателен"
```bash
cd voice_test
./setup_turn.sh
```

### Проблема: Порт занят
```bash
# Проверьте, что порты свободны
netstat -tlnp | grep :3478
netstat -tlnp | grep :8080

# Или используйте другие порты в .env
TURN_PORT=3479
PORT=8081
```

### Проблема: Соединение не устанавливается
```bash
# Для локального тестирования используйте 127.0.0.1
# Для удаленных клиентов нужен публичный IP
TURN_PUBLIC_IP=ВАШ_ПУБЛИЧНЫЙ_IP
```
```

## ⚙️ Как это работает

### Автоматическое P2P соединение

1. **TURN сервер** запускается автоматически вместе с SFU
2. **WebRTC клиенты** получают TURN credentials через сигнализацию
3. **ICE кандидаты** собираются через TURN для обхода NAT
4. **P2P соединение** устанавливается автоматически

### Архитектура

- **SFU (Selective Forwarding Unit)**: Маршрутизация медиа-потоков
- **TURN сервер (pion/turn)**: Обход NAT и firewall
- **WebSocket сигнализация**: Обмен SDP offer/answer
- **Docker**: Контейнеризация

## 🔧 Настройка

### Переменные окружения (.env)

```env
PORT=8080

# TURN сервер настройки (автоматически настраиваются)
TURN_PUBLIC_IP=ваш-публичный-ip
TURN_PORT=3478
TURN_USERNAME=ancord-user
TURN_PASSWORD=secure-password
```

### Ручная настройка IP

```bash
# Получить IP вручную
curl ifconfig.me

# Обновить .env
nano .env
# Изменить TURN_PUBLIC_IP=ваш-ip
```

## 🐳 Docker развертывание

### Docker:
```bash
docker build -t ancord-sfu:latest .
docker run -p 8080:8080 -p 3478:3478/udp --env-file .env ancord-sfu:latest
```

### Docker Compose:
```bash
docker-compose up -d
```

## 🔒 Безопасность

- TURN credentials генерируются динамически (24 часа)
- WebSocket с проверкой origin
- Переменные окружения для чувствительных данных
- Собственный TURN сервер (pion/turn)

## 📊 Мониторинг

Сервер логирует:
- Запуск TURN сервера
- WebRTC соединения
- Ошибки сигнализации
- RTP потоки

## 🛠️ Устранение неполадок

### Проблема: "TURN_PUBLIC_IP обязателен"
```bash
# Запустите настройку
./setup_turn.sh
```

### Проблема: Соединение не устанавливается
```bash
# Проверьте логи
docker logs ancord-sfu

# Проверьте порты
netstat -tlnp | grep :3478
netstat -tlnp | grep :8080
```

### Проблема: Docker не запускается
```bash
# Остановите предыдущие контейнеры
docker stop ancord-sfu
docker rm ancord-sfu

# Запустите заново
docker-compose up -d
```

# Запустите автоматическое развертывание
chmod +x deploy_wsl.sh
./deploy_wsl.sh
```

### 3. Настройка TURN сервера

```bash
# Получите публичный IP
./get_public_ip.sh

# Отредактируйте .env файл
nano .env
# Установите TURN_PUBLIC_IP=ВАШ_IP_АДРЕС
```

### 4. Запуск сервера

```bash
# Локально
go run ./sfu

# Или с Docker
docker-compose up -d
```

### 5. Тестирование

```bash
# В другом терминале запустите HTTP сервер для тестов
python3 -m http.server 8000
# Откройте http://localhost:8000 в браузере
```

## 🛠️ Управление проектом (Make)

Проект включает Makefile для удобного управления:

```bash
# Проверка системы
make check        # или ./check_wsl.sh

# Сборка и запуск
make build        # Собрать проект
make run          # Запустить сервер
make clean        # Очистить сборку

# Docker
make docker-build # Собрать образ
make docker-run   # Запустить контейнер
make docker-stop  # Остановить контейнер

# Docker Compose
make compose-up   # Запустить сервисы
make compose-down # Остановить сервисы

# Утилиты
make ip           # Получить публичный IP
make logs         # Посмотреть логи
make status       # Статус системы
```

## 🔧 Оптимизация WSL

Для лучшей производительности скопируйте `.wslconfig.example` в `C:\Users\<username>\.wslconfig`:

```powershell
# В PowerShell
copy-item .\voice_test\.wslconfig.example $env:USERPROFILE\.wslconfig
wsl --shutdown
wsl
```

## 🐧 Альтернативный запуск на WSL

### Ручная установка зависимостей

```bash
# Обновление пакетов
sudo apt update

# Установка Go
sudo apt install -y golang-go

# Установка curl (для получения IP)
sudo apt install -y curl

# Установка Docker (опционально)
sudo apt install -y docker.io
sudo systemctl start docker
sudo usermod -aG docker $USER
```

### Сборка и запуск

```bash
# Скачивание зависимостей
go mod download

# Сборка проекта
go build ./sfu

# Запуск
./sfu
```

## 🐳 Docker развертывание

### Docker:
```bash
cd voice_test
docker build -t ancord-sfu:latest .
docker run -p 8080:8080 -p 3478:3478/udp --env-file .env ancord-sfu:latest
```

### Docker Compose:
```bash
cd voice_test
docker-compose up -d
```

## ⚙️ Настройка TURN сервера

### Переменные окружения (.env)

```env
PORT=8080

# TURN сервер настройки (ОБЯЗАТЕЛЬНО)
TURN_PUBLIC_IP=your-server-public-ip
TURN_PORT=3478
TURN_USERNAME=ancord-user
TURN_PASSWORD=your-secure-password
```

### Настройка TURN_PUBLIC_IP

**ОБЯЗАТЕЛЬНО** укажите публичный IP адрес вашего сервера:

```bash
# Автоматически получить IP
./get_public_ip.sh

# Или вручную
curl ifconfig.me
# или
curl api.ipify.org
```

Без корректного TURN_PUBLIC_IP сервер не запустится!

## 🏗️ Архитектура

- **SFU**: Selective Forwarding Unit для эффективной маршрутизации медиа
- **TURN**: Собственный TURN сервер для обхода NAT и firewall
- **WebSocket**: Сигнализация для WebRTC
- **Docker**: Контейнеризация для легкого развертывания

## 🔒 Безопасность

- TURN credentials действительны 24 часа
- WebSocket с проверкой origin
- Переменные окружения для чувствительных данных
- Собственный TURN сервер (без зависимостей от внешних сервисов)
curl ifconfig.me
# или
curl api.ipify.org
```

Без корректного TURN_PUBLIC_IP сервер не запустится!

## Архитектура

- **SFU**: Selective Forwarding Unit для эффективной маршрутизации медиа
- **TURN**: Собственный TURN сервер для обхода NAT и firewall
- **WebSocket**: Сигнализация для WebRTC
- **Docker**: Контейнеризация для легкого развертывания

## Безопасность

- TURN credentials действительны 24 часа
- WebSocket с проверкой origin
- Переменные окружения для чувствительных данных
- Собственный TURN сервер (без зависимостей от внешних сервисов)