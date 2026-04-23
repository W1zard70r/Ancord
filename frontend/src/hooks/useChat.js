import { useEffect, useState, useRef } from 'react';

export const useChat = (token, chatID) => {
    const [messages, setMessages] = useState([]);
    const ws = useRef(null);

    useEffect(() => {
        if (!token || !chatID) return;

        // Закрываем старый сокет при смене чата
        if (ws.current) ws.current.close();
        // Очищаем сообщения
        setMessages([]);

        const socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);
        ws.current = socket;

        socket.onopen = () => {
            console.log("WS Connected");
            socket.send(JSON.stringify({ type: "join", chat_id: chatID }));
        };

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);

            if (data.type === "history") {
                setMessages(data.data || []);
            } else if (data.type === "message") {
                // Если это сообщение из БД, распарсим его content
                const msg = typeof data.message === 'string' ? JSON.parse(data.message) : data.message;
                setMessages((prev) => [...prev, msg]);
            } else if (data.type === "voice_status") {
                // Обрабатываем события из Kafka
                const actionText = data.data.action === 'join' ? 'подключился к голосовому каналу' : 'покинул голосовой канал';
                const sysMsg = {
                    isSystem: true, // Флаг для красивого отображения в App.jsx
                    content: `Пользователь ${data.data.user_id.split('-')[0]}... ${actionText}`
                };
                setMessages((prev) => [...prev, sysMsg]);
            }
        };

        return () => socket.close();
    }, [token, chatID]);

    const sendMessage = (content) => {
        if (!content.trim()) return;
        if (ws.current?.readyState === WebSocket.OPEN) {
            ws.current.send(JSON.stringify({ type: "message", chat_id: chatID, content }));
        }
    };

    return { messages, sendMessage };
};