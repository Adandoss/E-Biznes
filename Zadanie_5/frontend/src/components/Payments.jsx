import { useState, useEffect } from 'react';
import axios from 'axios';

import { useCart } from '../context/CartContext';

export default function Payments() {
  const { totalAmount, clearCart } = useCart();
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');
  const [payments, setPayments] = useState([]);

  const fetchPayments = () => {
    axios.get('http://localhost:8080/payments')
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
    axios.post('http://localhost:8080/payments', {
      amount: totalAmount,
      status: 'PAID'
    })
      .then(() => {
        setMessage('Zapłacono pomyślnie!');
        setLoading(false);
        clearCart();
        fetchPayments();
      })
      .catch(() => {
        setMessage('Błąd płatności.');
        setLoading(false);
      });
  };

  return (
    <div>
      <h2 className="text-lg font-bold mb-2">Płatności</h2>
      <p className="mb-2">Do zapłaty: {totalAmount.toFixed(2)} PLN</p>
      <button
        onClick={handleSubmit}
        disabled={loading || totalAmount <= 0}
        className="border px-2 py-1 bg-gray-200"
      >
        {loading ? 'Przetwarzanie...' : 'Kup'}
      </button>
      {message && <p className="mt-2 text-sm">{message}</p>}

      <div className="mt-8 border-t pt-4">
        <h3 className="font-bold mb-2">Historia płatności</h3>
        <ul className="list-disc pl-5">
          {payments && payments.length > 0 ? (
            payments.map(payment => (
              <li key={payment.ID} className="mb-1">
                Płatność nr {payment.ID} - {payment.amount.toFixed(2)} PLN ({payment.status})
              </li>
            ))
          ) : (
            <p>Brak historii płatności.</p>
          )}
        </ul>
      </div>
    </div>
  );
}
