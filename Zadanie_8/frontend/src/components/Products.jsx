import { useState, useEffect } from 'react';
import api from '../api';
import { useCart } from '../context/CartContext';

export default function Products() {
  const { addToCart } = useCart();
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);

  const [newName, setNewName] = useState('');
  const [newDesc, setNewDesc] = useState('');
  const [newPrice, setNewPrice] = useState('');
  const [newCategoryId, setNewCategoryId] = useState('');

  const fetchProducts = () => {
    api.get('/products')
      .then(response => {
        setProducts(response.data);
        setLoading(false);
      })
      .catch(error => {
        console.error('Błąd pobierania produktów:', error);
        setLoading(false);
      });
  };

  const fetchCategories = () => {
    api.get('/categories')
      .then(response => {
        setCategories(response.data);
        if (response.data && response.data.length > 0) {
          setNewCategoryId(response.data[0].ID);
        }
      })
      .catch(error => {
        console.error('Błąd pobierania kategorii:', error);
      });
  };

  useEffect(() => {
    fetchProducts();
    fetchCategories();
  }, []);

  const addProduct = (e) => {
    e.preventDefault();
    if (!newName || !newPrice || !newCategoryId) return;

    api.post('/products', {
      name: newName,
      description: newDesc,
      price: Number.parseFloat(newPrice),
      category_id: Number.parseInt(newCategoryId, 10),
    }).then(() => {
      setNewName('');
      setNewDesc('');
      setNewPrice('');
      fetchProducts();
    }).catch(err => console.error("Błąd dodawania produktu:", err));
  };

  return (
    <div className="space-y-8">
      <div>
        <h2 className="text-2xl font-bold text-gray-800 mb-6">Nasze Produkty</h2>
        {loading ? (
          <p className="text-gray-500 text-sm">Ładowanie produktów...</p>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {products?.map(product => (
              <div key={product.ID} className="bg-white p-6 rounded-lg shadow border border-gray-200 flex flex-col justify-between">
                <div>
                  <div className="flex justify-between items-start mb-2 gap-2">
                    <h3 className="text-lg font-bold text-gray-900">{product.name}</h3>
                    <span className="bg-blue-100 text-blue-800 text-xs px-2.5 py-1 rounded-full font-medium whitespace-nowrap">
                      {product.category?.name || 'Brak kategorii'}
                    </span>
                  </div>
                  <p className="text-gray-600 text-sm mb-4">{product.description}</p>
                </div>
                <div className="flex items-center justify-between mt-4">
                  <span className="text-xl font-extrabold text-blue-600">{product.price.toFixed(2)} PLN</span>
                  <button 
                    onClick={() => addToCart(product)} 
                    className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm font-semibold transition-colors shadow"
                  >
                    Dodaj do koszyka
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="bg-white p-6 rounded-lg shadow border border-gray-200 max-w-xl">
        <h3 className="text-xl font-bold text-gray-800 mb-4">Dodaj nowy produkt</h3>
        <form onSubmit={addProduct} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Nazwa produktu</label>
            <input
              type="text"
              placeholder="np. Słuchawki bezprzewodowe"
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Opis</label>
            <textarea
              placeholder="Opis produktu..."
              value={newDesc}
              onChange={(e) => setNewDesc(e.target.value)}
              className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              rows="2"
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Cena (PLN)</label>
              <input
                type="number"
                step="0.01"
                placeholder="0.00"
                value={newPrice}
                onChange={(e) => setNewPrice(e.target.value)}
                className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Kategoria</label>
              <select
                value={newCategoryId}
                onChange={(e) => setNewCategoryId(e.target.value)}
                className="w-full border border-gray-300 rounded px-3 py-2 text-sm bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
              >
                {categories.map(cat => (
                  <option key={cat.ID} value={cat.ID}>
                    {cat.name}
                  </option>
                ))}
              </select>
            </div>
          </div>
          <button 
            type="submit" 
            className="w-full bg-green-600 hover:bg-green-700 text-white py-2 rounded text-sm font-bold transition-colors shadow"
          >
            Dodaj produkt do sklepu
          </button>
        </form>
      </div>
    </div>
  );
}
