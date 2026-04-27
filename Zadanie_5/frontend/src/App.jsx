import { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, NavLink } from 'react-router-dom';
import axios from 'axios';
import { CartProvider } from './context/CartContext';
import Products from './components/Products';
import Cart from './components/Cart';
import Payments from './components/Payments';

function App() {
  const [cartId] = useState(1);

  return (
    <Router>
      <CartProvider cartId={cartId}>
        <div className="p-4">
          <header className="mb-4">
            <h1 className="text-xl font-bold mb-2">Sklep</h1>
            <nav className="flex gap-4 border-b pb-2">
              <NavLink to="/" className={({ isActive }) => isActive ? "font-bold underline" : ""}>Produkty</NavLink>
              <NavLink to="/cart" className={({ isActive }) => isActive ? "font-bold underline" : ""}>Koszyk</NavLink>
              <NavLink to="/payments" className={({ isActive }) => isActive ? "font-bold underline" : ""}>Płatności</NavLink>
            </nav>
          </header>

          <main>
            <Routes>
              <Route path="/" element={<Products />} />
              <Route path="/cart" element={<Cart />} />
              <Route path="/payments" element={<Payments />} />
            </Routes>
          </main>
        </div>
      </CartProvider>
    </Router>
  );
}

export default App;
