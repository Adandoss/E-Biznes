import { useState, useRef, useEffect } from 'react';
import axios from 'axios';

export default function Chatbot() {
  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState([
    {
      id: 1,
      sender: 'bot',
      text: 'Cześć! Jestem asystentem sklepu. Pomogę Ci w sprawach ubrań i naszego sklepu. O co chcesz zapytać?',
    }
  ]);
  const [inputText, setInputText] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const messagesEndRef = useRef(null);

  const KTOR_API_URL = 'http://localhost:8081/api/chat';

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSend = async (e) => {
    e.preventDefault();
    if (!inputText.trim()) return;

    const userText = inputText;
    setInputText('');

    const userMsg = { id: Date.now(), sender: 'user', text: userText };
    setMessages((prev) => [...prev, userMsg]);
    setIsLoading(true);

    try {
      //3.5 /////////////////////
      // Połączenie z kotlinem
      const response = await axios.post(KTOR_API_URL, { message: userText });
      const botMsg = { id: Date.now() + 1, sender: 'bot', text: response.data.response };
      setMessages((prev) => [...prev, botMsg]);
    } catch (error) {
      const errorMsg = { id: Date.now() + 2, sender: 'bot', text: 'Błąd połączenia z serwerem.' };
      setMessages((prev) => [...prev, errorMsg]);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="fixed bottom-5 right-5 z-50 font-sans">
      {!isOpen && (
        <button
          onClick={() => setIsOpen(true)}
          className="bg-blue-600 hover:bg-blue-700 text-white border-none rounded-full w-14 h-14 text-2xl cursor-pointer shadow-lg transition-colors flex items-center justify-center"
        >
          💬
        </button>
      )}

      {isOpen && (
        <div className="w-80 h-[400px] bg-white border border-gray-200 rounded-lg shadow-xl flex flex-col overflow-hidden">
          {/* Header */}
          <div className="bg-blue-600 text-white px-4 py-3 flex justify-between items-center">
            <span className="text-sm font-bold">Czat ze sklepem</span>
            <button
              onClick={() => setIsOpen(false)}
              className="bg-transparent text-white border-none text-base cursor-pointer hover:text-blue-200 transition-colors"
            >
              ✕
            </button>
          </div>

          {/* Messages */}
          <div className="flex-grow px-3 py-3 overflow-y-auto bg-gray-50 space-y-2">
            {messages.map((msg) => (
              <div
                key={msg.id}
                className={`flex ${msg.sender === 'user' ? 'justify-end' : 'justify-start'}`}
              >
                <div
                  className={`inline-block px-3 py-2 rounded-xl text-sm max-w-[80%] ${msg.sender === 'user'
                      ? 'bg-blue-600 text-white'
                      : 'bg-gray-200 text-gray-800'
                    }`}
                >
                  {msg.text}
                </div>
              </div>
            ))}
            {isLoading && (
              <p className="text-xs text-gray-400">Bot pisze...</p>
            )}
            <div ref={messagesEndRef} />
          </div>

          {/* Input */}
          <form onSubmit={handleSend} className="flex border-t border-gray-200">
            <input
              type="text"
              value={inputText}
              onChange={(e) => setInputText(e.target.value)}
              placeholder="Napisz coś..."
              disabled={isLoading}
              className="flex-grow px-3 py-2.5 border-none outline-none text-sm bg-white placeholder-gray-400"
            />
            <button
              type="submit"
              disabled={isLoading}
              className="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white border-none px-4 cursor-pointer text-sm font-medium transition-colors"
            >
              Wyślij
            </button>
          </form>
        </div>
      )}
    </div>
  );
}
