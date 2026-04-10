#!/bin/bash

# Читаем .env
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

BASE_URL="http://localhost:8080"

# Функция для вызова API
run_req() {
    local method=$1
    local path=$2
    local body=$3
    local token=$4
    
    local cmd="curl -s -w '\nHTTP_CODE:%{http_code}\n' -X $method $BASE_URL$path"
    if [ ! -z "$token" ]; then cmd="$cmd -H 'Authorization: Bearer $token'"; fi
    if [ ! -z "$body" ]; then cmd="$cmd -H 'Content-Type: application/json' -d '$body'"; fi
    eval "$cmd"
}

# 1. Регистрация
echo "1. Регистрация юзеров..."
U1_JSON=$(run_req "POST" "/register" '{"username": "user1", "password": "password123"}')
U2_JSON=$(run_req "POST" "/register" '{"username": "user2", "password": "password456"}')
U3_JSON=$(run_req "POST" "/register" '{"username": "user3", "password": "password789"}')

U1_ID=$(echo "$U1_JSON" | head -n1 | jq -r .id)
U2_ID=$(echo "$U2_JSON" | head -n1 | jq -r .id)
U3_ID=$(echo "$U3_JSON" | head -n1 | jq -r .id)

# 2. Логин (получаем токены)
echo "2. Логин..."
T1=$(run_req "POST" "/login" '{"username": "user1", "password": "password123"}' | head -n1 | jq -r .token)
T2=$(run_req "POST" "/login" '{"username": "user2", "password": "password456"}' | head -n1 | jq -r .token)

# 3. Создание чатов
echo "3. Создание чатов..."
# Общий чат
C_GEN=$(run_req "POST" "/chats" '{"name": "Общий чат"}' "$T1" | head -n1 | jq -r .id)
# Приватные (user1-user2, user2-user3, user1-user3)
C_12=$(run_req "POST" "/chats" '{"name": "User1 & User2"}' "$T1" | head -n1 | jq -r .id)
C_23=$(run_req "POST" "/chats" '{"name": "User2 & User3"}' "$T2" | head -n1 | jq -r .id)
C_13=$(run_req "POST" "/chats" '{"name": "User1 & User3"}' "$T1" | head -n1 | jq -r .id)
# Пустой чат
C_EMPTY=$(run_req "POST" "/chats" '{"name": "Пустой чат"}' "$T1" | head -n1 | jq -r .id)

# 4. Добавление участников
echo "4. Добавление участников..."
# Общий чат (добавляем всех)
run_req "POST" "/chats/$C_GEN/members" "{\"target_user_id\": \"$U2_ID\"}" "$T1" > /dev/null
run_req "POST" "/chats/$C_GEN/members" "{\"target_user_id\": \"$U3_ID\"}" "$T1" > /dev/null
# Приватные
run_req "POST" "/chats/$C_12/members" "{\"target_user_id\": \"$U2_ID\"}" "$T1" > /dev/null
run_req "POST" "/chats/$C_23/members" "{\"target_user_id\": \"$U3_ID\"}" "$T2" > /dev/null
run_req "POST" "/chats/$C_13/members" "{\"target_user_id\": \"$U3_ID\"}" "$T1" > /dev/null

# 5. Сообщения
echo "5. Отправка сообщений..."
run_req "POST" "/message" "{\"chat_id\": \"$C_GEN\", \"content\": \"Привет, это общий чат!\"}" "$T1" > /dev/null
run_req "POST" "/message" "{\"chat_id\": \"$C_12\", \"content\": \"Секретик для User2\"}" "$T1" > /dev/null
run_req "POST" "/message" "{\"chat_id\": \"$C_23\", \"content\": \"Секретик для User3\"}" "$T2" > /dev/null

echo "--- Seed completed ---"