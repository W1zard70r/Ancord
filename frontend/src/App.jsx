import { useState, useRef, useEffect } from 'react';
import { useChat } from './hooks/useChat';
import axios from 'axios';

function App() {
  const [username, setUsername] = useState('user1');
  const [password, setPassword] = useState('password123');
  const [token, setToken] = useState('');

  // ВАЖНО: Убрали хардкод!
  const [chatIDInput, setChatIDInput] = useState('');
  const [activeChatID, setActiveChatID] = useState('');

  const [input, setInput] = useState('');
  const { messages, sendMessage, connect } = useChat(token, activeChatID);
  const messagesEndRef = useRef(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  // Если activeChatID изменился и есть токен - коннектимся автоматически
  useEffect(() => {
    if (token && activeChatID) {
      connect();
    }
  }, [activeChatID, token, connect]);

  const handleLogin = async () => {
    try {
      // Для Docker используем относительный путь или ENV.
      // Пока оставим localhost, так как браузер друга будет стучаться на свой же localhost, 
      // если он поднимет весь Docker у себя.
      const res = await axios.post('http://localhost:8080/login', { username, password });
      setToken(res.data.token);
    } catch (e) { alert("Логин не удался"); }
  };

  const handleJoinChat = () => {
    if (!chatIDInput.trim()) return;
    setActiveChatID(chatIDInput.trim());
  };

  const handleSend = () => {
    if (input.trim() === '') return;
    sendMessage(input);
    setInput('');
  };

  return (
    <div className="h-screen bg-gray-100 flex items-center justify-center p-4">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-2xl overflow-hidden flex flex-col h-[80vh]">
        <div className="bg-blue-600 text-white p-4 flex justify-between items-center">
          <h1 className="text-xl font-bold">Ancord Chat</h1>
          {token && <div className="text-xs opacity-75">Logged in</div>}
        </div>

        <div className="flex-grow flex flex-col p-4 overflow-hidden">
          {!token ? (
            <div className="flex flex-col justify-center h-full space-y-4 max-w-xs mx-auto w-full">
              <h2 className="text-lg font-semibold text-center">Login</h2>
              <input className="border rounded p-2" placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} />
              <input className="border rounded p-2" type="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} />
              <button onClick={handleLogin} className="bg-blue-600 text-white rounded p-2">Sign In</button>
            </div>
          ) : !activeChatID ? (
            <div className="flex flex-col justify-center h-full space-y-4 max-w-md mx-auto w-full">
              <h2 className="text-lg font-semibold text-center">Join a Chat</h2>
              <input
                className="border rounded p-2 w-full font-mono text-sm"
                placeholder="Enter Chat UUID (e.g. 019d...)"
                value={chatIDInput}
                onChange={e => setChatIDInput(e.target.value)}
              />
              <button onClick={handleJoinChat} className="bg-green-500 text-white rounded p-2">Join Room</button>
              <p className="text-xs text-gray-400 text-center mt-4">
                (Создай чат через Postman/cURL и вставь его ID сюда)
              </p>
            </div>
          ) : (
            <div className="flex flex-col h-full">
              <div className="text-xs text-gray-500 mb-2">Room: <span className="font-mono">{activeChatID}</span></div>
              <div className="flex-grow border rounded bg-gray-50 p-4 overflow-y-auto mb-4 space-y-3">
                {messages.map((m, i) => (
                  <div key={m.id || i} className="bg-white p-3 rounded shadow-sm border border-gray-100">
                    <div className="flex justify-between items-baseline mb-1">
                      <span className="text-xs font-bold text-blue-600">User: {m.user_id?.split('-')[0]}</span>
                    </div>
                    <p className="text-gray-800">{m.content}</p>
                  </div>
                ))}
                <div ref={messagesEndRef} />
              </div>
              <div className="flex gap-2 flex-shrink-0">
                <input className="border rounded p-3 flex-grow" value={input} onChange={(e) => setInput(e.target.value)} onKeyPress={(e) => e.key === 'Enter' && handleSend()} placeholder="Message..." />
                <button onClick={handleSend} className="bg-blue-600 text-white rounded px-6 py-2">Send</button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
export default App;