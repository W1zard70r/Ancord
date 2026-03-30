import { useState, useRef, useEffect } from 'react';
import { useChat } from './hooks/useChat';
import axios from 'axios';

function App() {
  const [username, setUsername] = useState('user1');
  const [password, setPassword] = useState('password123');
  const [token, setToken] = useState('');

  // Вставь сюда реальный ID чата из своей базы
  const [chatID, setChatID] = useState('019d3b68-0b24-732d-827e-3b41aa37cf22');
  const [input, setInput] = useState('');

  const { messages, sendMessage, connect } = useChat(token, chatID);
  const messagesEndRef = useRef(null);

  // Автоскролл вниз при новых сообщениях
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const handleLogin = async () => {
    try {
      const res = await axios.post('http://localhost:8080/login', {
        username,
        password
      });
      setToken(res.data.token);
    } catch (e) {
      console.error(e);
      alert("Логин не удался. Проверь консоль.");
    }
  };

  const handleSend = () => {
    if (input.trim() === '') return;
    sendMessage(input);
    setInput('');
  };

  return (
    <div className="h-screen bg-gray-100 flex items-center justify-center p-4">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-2xl overflow-hidden flex flex-col h-[80vh]">

        {/* Хедер */}
        <div className="bg-blue-600 text-white p-4">
          <h1 className="text-xl font-bold">Ancord Chat</h1>
          {token && <div className="text-xs opacity-75 mt-1">Logged in</div>}
        </div>

        {/* Тело (Авторизация ИЛИ Чат) */}
        <div className="flex-grow flex flex-col p-4 overflow-hidden">

          {!token ? (
            <div className="flex flex-col justify-center h-full space-y-4 max-w-xs mx-auto w-full">
              <h2 className="text-lg font-semibold text-center">Login to Chat</h2>
              <input
                className="border rounded p-2 focus:ring-2 focus:ring-blue-500 outline-none"
                placeholder="Username"
                value={username}
                onChange={e => setUsername(e.target.value)}
              />
              <input
                className="border rounded p-2 focus:ring-2 focus:ring-blue-500 outline-none"
                type="password"
                placeholder="Password"
                value={password}
                onChange={e => setPassword(e.target.value)}
              />
              <button onClick={handleLogin} className="bg-blue-600 text-white rounded p-2 font-medium hover:bg-blue-700 transition">
                Sign In
              </button>
            </div>
          ) : (
            <div className="flex flex-col h-full">

              {/* Кнопка коннекта (пока оставим вручную для надежности) */}
              <button
                onClick={() => connect()}
                className="bg-green-500 text-white rounded p-2 mb-4 font-medium hover:bg-green-600 transition flex-shrink-0"
              >
                Connect to Room
              </button>

              {/* Список сообщений */}
              <div className="flex-grow border rounded bg-gray-50 p-4 overflow-y-auto mb-4 space-y-3">
                {messages.length === 0 ? (
                  <div className="text-center text-gray-400 mt-10">No messages yet...</div>
                ) : (
                  messages.map((m, i) => (
                    <div key={m.id || i} className="bg-white p-3 rounded shadow-sm border border-gray-100">
                      <div className="flex justify-between items-baseline mb-1">
                        <span className="text-xs font-bold text-blue-600" title={m.user_id}>User ID: {m.user_id.split('-')[0]}...</span>
                        {m.created_at && (
                          <span className="text-[10px] text-gray-400">
                            {new Date(m.created_at).toLocaleTimeString()}
                          </span>
                        )}
                      </div>
                      <p className="text-gray-800 break-words">{m.content}</p>
                    </div>
                  ))
                )}
                <div ref={messagesEndRef} /> {/* Якорь для автоскролла */}
              </div>

              {/* Ввод сообщения */}
              <div className="flex gap-2 flex-shrink-0">
                <input
                  className="border rounded p-3 flex-grow focus:ring-2 focus:ring-blue-500 outline-none"
                  value={input}
                  onChange={(e) => setInput(e.target.value)}
                  onKeyPress={(e) => e.key === 'Enter' && handleSend()}
                  placeholder="Type a message..."
                />
                <button
                  onClick={handleSend}
                  className="bg-blue-600 text-white rounded px-6 py-2 font-medium hover:bg-blue-700 transition"
                >
                  Send
                </button>
              </div>

            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default App;