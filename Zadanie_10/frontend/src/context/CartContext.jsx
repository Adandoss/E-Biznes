import { createContext, useState, useContext, useEffect, useMemo, useCallback, useRef } from 'react';
import api from '../api';

const CartContext = createContext();

export const useCart = () => useContext(CartContext);

export const CartProvider = ({ children }) => {
  const [cartItems, setCartItems] = useState([]);
  const [cartId, setCartId] = useState(null);
  const [loading, setLoading] = useState(true);
  const cartItemsRef = useRef(cartItems);

  useEffect(() => {
    cartItemsRef.current = cartItems;
  }, [cartItems]);

  useEffect(() => {
    api.get('/carts/mine')
      .then(res => {
        setCartId(res.data.ID);
        if (res.data?.items) {
          setCartItems(res.data.items);
        }
        setLoading(false);
      })
      .catch(err => {
        console.error("Error fetching cart:", err);
        setLoading(false);
      });
  }, []);

  const addToCart = useCallback(async (product) => {
    if (!cartId) return;
    try {
      const response = await api.post(`/carts/${cartId}/items`, {
        product_id: product.ID,
        quantity: 1
      });
      setCartItems(prev => [...prev, response.data]);
    } catch (error) {
      console.error("Error adding to cart:", error);
    }
  }, [cartId]);

  const removeFromCart = useCallback(async (itemId) => {
    if (!cartId) return;
    const previousItems = cartItemsRef.current;
    setCartItems(prev => prev.filter(item => Number(item.ID) !== Number(itemId)));
    try {
      await api.delete(`/carts/${cartId}/items/${itemId}`);
    } catch (error) {
      console.error("Error removing from cart:", error);
      setCartItems(previousItems);
    }
  }, [cartId]);

  const clearCart = useCallback(async () => {
    if (!cartId) return;
    try {
      await api.delete(`/carts/${cartId}`);
      setCartItems([]);
      const res = await api.get('/carts/mine');
      setCartId(res.data.ID);
    } catch (error) {
      console.error("Error clearing cart:", error);
    }
  }, [cartId]);

  const totalAmount = useMemo(() =>
    cartItems.reduce((sum, item) => sum + ((item.product?.price || 0) * (item.quantity || 0)), 0),
    [cartItems]
  );

  const value = useMemo(() => ({
    cartItems,
    addToCart,
    removeFromCart,
    clearCart,
    totalAmount,
    loading
  }), [cartItems, addToCart, removeFromCart, clearCart, totalAmount, loading]);

  return (
    <CartContext.Provider value={value}>
      {children}
    </CartContext.Provider>
  );
};
