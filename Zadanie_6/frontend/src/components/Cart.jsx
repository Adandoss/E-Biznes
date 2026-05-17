import { useState, useEffect } from 'react';
import axios from 'axios';

import { useCart } from '../context/CartContext';

export default function Cart() {
  const { cartItems, removeFromCart, totalAmount, loading } = useCart();

  return (
    <div>
      <h2 className="text-lg font-bold mb-2">Koszyk</h2>
      {loading ? <p>Ładowanie...</p> : (
        <>
          {cartItems && cartItems.length > 0 ? (
            <ul className="list-disc pl-5">
              {cartItems.map(item => (
                <li key={item.ID} className="mb-2">
                  {item.product.name} ({item.quantity} szt.) - {item.product.price * item.quantity} PLN
                  <button onClick={() => removeFromCart(item.ID)} className="ml-2 border px-1 bg-gray-100 text-sm">Usuń</button>
                </li>
              ))}
            </ul>
          ) : (
            <p>Koszyk jest pusty.</p>
          )}
          <p className="mt-4 font-bold">Suma: {totalAmount.toFixed(2)} PLN</p>
        </>
      )}
    </div>
  );
}
