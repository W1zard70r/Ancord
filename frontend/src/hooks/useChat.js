import { useEffect, useState, useCallback, useRef } from 'react';

export const useChat = (token, chatID) => {
    const [messages, setMessages] = useState([]);
    const ws = useRef(null); // Используем useRef, чтобы не терять сокет при ререндерах

    const connect = useCallback(() => {
        if (!token || !chatID) return;

        // Если сокет уже открыт - закрываем старый
        if (ws.current) {
            ws.current.close();
        }

        const socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

        socket.onopen = () => {
            console.log("WebSocket Connected");
            // Отправляем join, чтобы сервер подписал нас и прислал историю
            socket.send(JSON.stringify({ type: "join", chat_id: chatID }));
        };

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            console.log("Received from server:", data);

            if (data.type === "history") {
                // История приходит в массиве data.data
                setMessages(data.data || []);
            } else if (data.type === "message") {
                // Новое сообщение приходит в объекте data.message (мы поменяли это на бэке!)
                if (data.message) {
                    setMessages((prev) => [...prev, data.message]);
                }
            } else if (data.type === "error") {
                console.error("Server Error:", data.message);
                alert("Ошибка: " + data.message);
            }
        };

        socket.onclose = () => console.log("WebSocket Disconnected");
        socket.onerror = (error) => console.error("WebSocket Error:", error);

        ws.current = socket;
    }, [token, chatID]);

    const sendMessage = useCallback((content) => {
        if (ws.current && ws.current.readyState === WebSocket.OPEN) {
            ws.current.send(JSON.stringify({
                type: "message",
                chat_id: chatID,
                content: content
            }));
        } else {
            console.error("Cannot send message: WebSocket is not open");
        }
    }, [chatID]);

    // Закрываем сокет при размонтировании компонента
    useEffect(() => {
        return () => {
            if (ws.current) {
                ws.current.close();
            }
        };
    }, []);

    return { messages, sendMessage, connect };
};