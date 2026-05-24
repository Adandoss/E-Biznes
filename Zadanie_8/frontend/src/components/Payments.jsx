import { useState, useEffect } from 'react';
import api from '../api';
import { useCart } from '../context/CartContext';

export default function Payments() {
  const { totalAmount, clearCart } = useCart();
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');
  const [payments, setPayments] = useState([]);

  const fetchPayments = () => {
    api.get('/payments')
      .then(res => setPayments(res.data))
      .catch(err => console.error("Błąd pobierania płatności", err));
  };

  useEffect(() => {
    fetchPayments();
  }, []);

  const handleSubmit = () => {
    if (totalAmount <= 0) {
      setMessage('Koszyk jest pusty.');
      return;
    }

    setLoading(true);
    api.post('/payments', {
      amount: totalAmount,
      status: 'PAID'
    })
      .then(() => {
        setMessage('Zapłacono pomyślnie! Zamówienie zostało zrealizowane.');
        setLoading(false);
        clearCart();
        fetchPayments();
      })
      .catch(() => {
        setMessage('Błąd podczas procesowania płatności.');
        setLoading(false);
      });
  };

  return (
    <div className="max-w-2xl mx-auto space-y-8">
      <div className="bg-white p-6 rounded-lg shadow border border-gray-200">
        <h2 className="text-2xl font-bold text-gray-800 mb-6">Płatność za zamówienie</h2>
        
        <div className="bg-blue-50 border border-blue-100 p-4 rounded mb-6 flex justify-between items-center">
          <span className="font-semibold text-blue-800">Kwota do zapłaty:</span>
          <span className="text-2xl font-extrabold text-blue-600">{totalAmount.toFixed(2)} PLN</span>
        </div>

        <button
          onClick={handleSubmit}
          disabled={loading || totalAmount <= 0}
          className="w-full bg-green-600 hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed text-white font-bold py-2.5 rounded shadow transition-colors"
        >
          {loading ? 'Przetwarzanie płatności...' : 'Zapłać teraz'}
        </button>

        {message && (
          <div className={`mt-4 p-3 rounded text-sm font-medium ${message.includes('pomyślnie') ? 'bg-green-50 text-green-800 border border-green-100' : 'bg-red-50 text-red-800 border border-red-100'}`}>
            {message}
          </div>
        )}
      </div>

      <div className="bg-white p-6 rounded-lg shadow border border-gray-200">
        <h3 className="text-xl font-bold text-gray-800 mb-4">Historia Płatności</h3>
        {payments && payments.length > 0 ? (
          <div className="divide-y divide-gray-100">
            {payments.map(payment => (
              <div key={payment.ID} className="py-3 flex justify-between items-center text-sm">
                <div>
                  <span className="font-semibold text-gray-800">Płatność #{payment.ID}</span>
                  <p className="text-xs text-gray-400">{new Date(payment.CreatedAt).toLocaleString()}</p>
                </div>
                <div className="flex items-center gap-3">
                  <span className="font-bold text-gray-900">{payment.amount.toFixed(2)} PLN</span>
                  <span className="bg-green-100 text-green-800 text-xs px-2 py-0.5 rounded font-medium">
                    {payment.status}
                  </span>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-gray-500 text-sm">Brak historii płatności.</p>
        )}
      </div>
    </div>
  );
}
