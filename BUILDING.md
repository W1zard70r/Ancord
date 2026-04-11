# Запуск
```bash
cp .env.example .env // создание .env файла
sudo make up //запускает контейнеры
sudo make migrate-up //создаёт таблицы в бд
chmod +x scripts/seed.sh // не уверен что нужна, но на всякий сделай (права на curl из .sh файла)
sudo make seed //делает пользователей и чаты
```
# Перезапустить (меняется код внутри, либы не переустанавливаются)
```bash
sudo make up
```

# Снести бд
```bash
sudo make migrate-down
```

# Снести контейнер
```bash
sudo make down
```