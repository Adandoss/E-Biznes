import { useCart } from '../context/CartContext';
import { Link } from 'react-router-dom';

export default function Cart() {
  const { cartItems, removeFromCart, totalAmount, loading } = useCart();

  return (
    <div className="max-w-2xl bg-white p-6 rounded-lg shadow border border-gray-200 mx-auto">
      <h2 className="text-2xl font-bold text-gray-800 mb-6">Twój Koszyk</h2>
      {loading ? (
        <p className="text-gray-500 text-sm">Ładowanie koszyka...</p>
      ) : (
        <>
          {cartItems && cartItems.length > 0 ? (
            <div className="space-y-4">
              {cartItems.map(item => (
                <div key={item.ID} className="flex justify-between items-center p-4 border border-gray-100 rounded bg-gray-50">
                  <div>
                    <h4 className="font-bold text-gray-900">{item.product?.name}</h4>
                    <p className="text-sm text-gray-500">{item.quantity} szt. x {item.product?.price?.toFixed(2)} PLN</p>
                  </div>
                  <div className="flex items-center gap-4">
                    <span className="font-semibold text-gray-900">{(item.product?.price * item.quantity).toFixed(2)} PLN</span>
                    <button 
                      onClick={() => removeFromCart(item.ID)} 
                      className="text-red-500 hover:text-red-700 hover:bg-red-50 px-3 py-1.5 rounded text-sm font-semibold transition-colors border border-red-200"
                    >
                      Usuń
                    </button>
                  </div>
                </div>
              ))}
              
              <div className="border-t pt-4 mt-6 flex justify-between items-center">
                <span className="text-lg font-medium text-gray-700">Suma do zapłaty:</span>
                <span className="text-2xl font-extrabold text-blue-600">{totalAmount.toFixed(2)} PLN</span>
              </div>

              <div className="mt-6">
                <Link 
                  to="/payments" 
                  className="block text-center w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2.5 rounded shadow transition-colors"
                >
                  Przejdź do płatności
                </Link>
              </div>
            </div>
          ) : (
            <div className="text-center py-8">
              <p className="text-gray-500 mb-4">Twój koszyk jest pusty.</p>
              <Link to="/" className="text-blue-600 hover:underline font-medium text-sm">
                Wróć do listy produktów i dodaj coś do koszyka!
              </Link>
            </div>
          )}
        </>
      )}
    </div>
  );
}
