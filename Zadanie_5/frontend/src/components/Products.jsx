import { useState, useEffect } from 'react';
import axios from 'axios';

import { useCart } from '../context/CartContext';

export default function Products() {
  const { addToCart } = useCart();
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);

  const [newName, setNewName] = useState('');
  const [newDesc, setNewDesc] = useState('');
  const [newPrice, setNewPrice] = useState('');

  const fetchProducts = () => {
    setLoading(true);
    axios.get('http://localhost:8080/products')
      .then(response => {
        setProducts(response.data);
        setLoading(false);
      })
      .catch(error => {
        console.error('Błąd:', error);
        setLoading(false);
      });
  };

  useEffect(() => {
    fetchProducts();
  }, []);


  const addProduct = (e) => {
    e.preventDefault();
    if (!newName || !newPrice) return;

    axios.post('http://localhost:8080/products', {
      name: newName,
      description: newDesc,
      price: parseFloat(newPrice),
    }).then(() => {
      setNewName('');
      setNewDesc('');
      setNewPrice('');
      fetchProducts();
    }).catch(err => console.error("Błąd dodawania produktu:", err));
  };

  return (
    <div>
      <h2 className="text-lg font-bold mb-2">Produkty</h2>
      {loading ? <p>Ładowanie...</p> : (
        <ul className="list-disc pl-5 mb-6">
          {products && products.map(product => (
            <li key={product.ID} className="mb-2">
              {product.name} - {product.price} PLN
              <button onClick={() => addToCart(product)} className="ml-2 border px-1 bg-gray-100 text-sm">
                Dodaj do koszyka
              </button>
            </li>
          ))}
        </ul>
      )}

      <div className="mt-8 border-t pt-4">
        <h3 className="font-bold mb-2">Dodaj nowy produkt</h3>
        <form onSubmit={addProduct} className="flex flex-col gap-2 max-w-sm">
          <input
            type="text"
            placeholder="Nazwa produktu"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            className="border p-1"
            required
          />
          <input
            type="text"
            placeholder="Opis"
            value={newDesc}
            onChange={(e) => setNewDesc(e.target.value)}
            className="border p-1"
          />
          <input
            type="number"
            step="0.01"
            placeholder="Cena (PLN)"
            value={newPrice}
            onChange={(e) => setNewPrice(e.target.value)}
            className="border p-1"
            required
          />
          <button type="submit" className="border px-2 py-1 bg-gray-200 w-max mt-1">
            Dodaj produkt
          </button>
        </form>
      </div>
    </div>
  );
}
