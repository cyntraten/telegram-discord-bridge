
# Telegram Discord Bridge

Мост для автоматической пересылки постов из Telegram-канала в Discord-канал.

## Возможности

- **Поддержка медиа**: пересылка текста, картинок, видео, гифок и альбомов (медиагрупп).
- **Безопасность**: фильтрация сообщений по ID целевого Telegram-канала.
- **Docker-ready**: быстрая сборка и запуск в изолированном контейнере.

---

## Быстрый старт

1. Клонируйте репозиторий:
```bash
git clone https://github.com/cyntraten/telegram-discord-bridge
cd telegram-discord-bridge
```


2. Создайте файл `.env` на основе примера (`.env.example`) и заполните его:
```env
TelegramToken="YOUR_TG_BOT_TOKEN"
DiscordToken="YOUR_DISCORD_BOT_TOKEN"
DiscordChannelId="YOUR_DISCORD_CHANNEL_ID"
TargetTelegramChannelId="-100XXXXXXXXXX"
```


> *Примечание: `TargetTelegramChannelId` используется для безопасности, чтобы бот игнорировал посторонние чаты.*


3. Соберите и запустите Docker-контейнер:
```bash
docker compose up --build -d
```


4. Просмотр логов приложения:
```bash
docker logs -f tg-discord-bridge
```



---

##Планы на будущее

* Интеграция с **FFmpeg** для автоматического сжатия медиафайлов под лимиты Discord и конвертации стикеров в поддерживаемые форматы.
