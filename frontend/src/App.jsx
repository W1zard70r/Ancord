import { useState } from 'react';
import { useChat } from './hooks/useChat';
import axios from 'axios';

function App() {
  const [auth, setAuth] = useState({ user: '', pass: '', token: '' });
  const [chats, setChats] = useState([]);
  const [activeChat, setActiveChat] = useState(null);
  const [input, setInput] = useState('');

  const { messages, sendMessage } = useChat(auth.token, activeChat?.id);

  const handleLogin = async () => {
    try {
      const res = await axios.post('http://localhost:8080/login', { username: auth.user, password: auth.pass });
      setAuth(prev => ({ ...prev, token: res.data.token }));

      // Загружаем список чатов
      const chatRes = await axios.get('http://localhost:8080/chats/my', {
        headers: { Authorization: `Bearer ${res.data.token}` }
      });
      setChats(chatRes.data || []);
    } catch (e) { alert("Ошибка авторизации"); }
  };

  return (
    <div className="flex h-screen bg-gray-100">
      <div className="w-64 bg-gray-900 text-white p-4">
        <h2 className="font-bold mb-4">Ancord</h2>
        {chats.map(c => (
          <button key={c.id} onClick={() => setActiveChat(c)} className="block w-full text-left p-2 hover:bg-gray-700 rounded">
            {c.name}
          </button>
        ))}
      </div>

      <div className="flex-grow flex flex-col bg-white">
        {!auth.token ? (
          <div className="m-auto w-64 space-y-2">
            <input className="border p-2 w-full" placeholder="User" onChange={e => setAuth({ ...auth, user: e.target.value })} />
            <input className="border p-2 w-full" type="password" placeholder="Pass" onChange={e => setAuth({ ...auth, pass: e.target.value })} />
            <button onClick={handleLogin} className="w-full bg-blue-600 text-white p-2">Login</button>
          </div>
        ) : activeChat ? (
          <>
            <div className="p-4 border-b font-bold">{activeChat.name}</div>
            <div className="flex-grow overflow-y-auto p-4 space-y-2">
              {messages.map((m, i) => <div key={i} className="p-2 bg-blue-100 rounded self-start">{m.content}</div>)}
            </div>
            <div className="p-4 border-t flex gap-2">
              <input className="border p-2 flex-grow" value={input} onChange={e => setInput(e.target.value)} />
              <button onClick={() => { sendMessage(input); setInput(''); }} className="bg-blue-600 text-white px-4 py-2 rounded">Send</button>
            </div>
          </>
        ) : <div className="m-auto text-gray-500">Выберите чат</div>}
      </div>
    </div>
  );
}
export default App;