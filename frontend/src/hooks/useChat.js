import { useEffect, useState, useRef } from 'react';

export const useChat = (token, chatID) => {
    const [messages, setMessages] = useState([]);
    const ws = useRef(null);

    useEffect(() => {
        if (!token || !chatID) return;

        // Закрываем старый сокет при смене чата
        if (ws.current) ws.current.close();

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
            }
        };

        return () => socket.close();
    }, [token, chatID]);

    const sendMessage = (content) => {
        if (ws.current?.readyState === WebSocket.OPEN) {
            ws.current.send(JSON.stringify({ type: "message", chat_id: chatID, content }));
        }
    };

    return { messages, sendMessage };
};