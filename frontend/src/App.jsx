import { useState, useEffect, useRef } from 'react';
import { useChat } from './hooks/useChat';
import { useVoice } from './hooks/useVoice';
import axios from 'axios';

// Невидимый плеер для воспроизведения чужих голосов
const AudioPlayer = ({ stream }) => {
  const audioRef = useRef(null);
  useEffect(() => {
    if (audioRef.current && stream) {
      audioRef.current.srcObject = stream;
    }
  }, [stream]);
  return <audio ref={audioRef} autoPlay className="hidden" />;
};

function App() {
  const [auth, setAuth] = useState({ user: '', pass: '', token: '', userId: '' });
  const [chats, setChats] = useState([]);
  const [activeChat, setActiveChat] = useState(null);
  const [input, setInput] = useState('');

  const { messages, sendMessage } = useChat(auth.token, activeChat?.id);
  const { status: voiceStatus, joinVoice, leaveVoice, remoteStreams } = useVoice(auth.userId, activeChat?.id);

  const handleLogin = async () => {
    try {
      const res = await axios.post('http://localhost:8080/login', { username: auth.user, password: auth.pass });
      const token = res.data.token;

      // Достаем user_id из JWT-токена
      const payload = JSON.parse(atob(token.split('.')[1]));

      setAuth(prev => ({ ...prev, token: token, userId: payload.user_id }));

      // Загружаем список чатов
      const chatRes = await axios.get('http://localhost:8080/chats/my', {
        headers: { Authorization: `Bearer ${token}` }
      });
      setChats(chatRes.data || []);
    } catch (e) { alert("Ошибка авторизации"); }
  };

  // Автоматически выходим из голосового, если переключили чат
  useEffect(() => {
    leaveVoice();
  }, [activeChat, leaveVoice]);

  return (
    <div className="flex h-screen bg-gray-100">
      <div className="w-64 bg-gray-900 text-white p-4">
        <h2 className="font-bold mb-4">Ancord</h2>
        {chats.map(c => (
          <button
            key={c.id}
            onClick={() => setActiveChat(c)}
            className={`block w-full text-left p-2 mb-1 rounded transition-colors ${activeChat?.id === c.id ? 'bg-blue-600' : 'hover:bg-gray-700'}`}
          >
            {c.name}
          </button>
        ))}
      </div>

      <div className="flex-grow flex flex-col bg-white">
        {!auth.token ? (
          <div className="m-auto w-64 space-y-2">
            <input className="border p-2 w-full rounded" placeholder="User" onChange={e => setAuth({ ...auth, user: e.target.value })} />
            <input className="border p-2 w-full rounded" type="password" placeholder="Pass" onChange={e => setAuth({ ...auth, pass: e.target.value })} />
            <button onClick={handleLogin} className="w-full bg-blue-600 hover:bg-blue-700 transition text-white p-2 rounded">Login</button>
          </div>
        ) : activeChat ? (
          <>
            <div className="p-4 border-b flex justify-between items-center bg-gray-50">
              <span className="font-bold text-lg">{activeChat.name}</span>

              <div className="flex items-center gap-4">
                {/* Рендерим скрытые плееры для каждого человека в ГС */}
                {remoteStreams.map(stream => <AudioPlayer key={stream.id} stream={stream} />)}

                <span className={`text-sm font-medium ${voiceStatus === 'connected' ? 'text-green-600' : 'text-gray-500'}`}>
                  Голос: {voiceStatus}
                </span>

                {voiceStatus === 'disconnected' || voiceStatus === 'error' ? (
                  <button onClick={joinVoice} className="bg-green-500 hover:bg-green-600 transition text-white px-4 py-2 rounded shadow">
                    Войти в ГС
                  </button>
                ) : (
                  <button onClick={leaveVoice} className="bg-red-500 hover:bg-red-600 transition text-white px-4 py-2 rounded shadow">
                    Выйти из ГС
                  </button>
                )}
              </div>
            </div>

            <div className="flex-grow overflow-y-auto p-4 space-y-2">
              {messages.map((m, i) => (
                m.isSystem ? (
                  <div key={i} className="text-center text-gray-400 text-sm italic my-2">{m.content}</div>
                ) : (
                  <div key={i} className="p-2 bg-blue-100 rounded self-start w-fit">{m.content}</div>
                )
              ))}
            </div>

            <div className="p-4 border-t flex gap-2">
              <input
                className="border p-2 flex-grow rounded"
                value={input}
                onChange={e => setInput(e.target.value)}
                onKeyDown={e => {
                  if (e.key === 'Enter') {
                    sendMessage(input);
                    setInput('');
                  }
                }}
                placeholder="Введите сообщение..."
              />
              <button
                onClick={() => { sendMessage(input); setInput(''); }}
                className="bg-blue-600 hover:bg-blue-700 transition text-white px-4 py-2 rounded"
              >
                Send
              </button>
            </div>
          </>
        ) : <div className="m-auto text-gray-500 text-lg">Выберите чат</div>}
      </div>
    </div>
  );
}
export default App;