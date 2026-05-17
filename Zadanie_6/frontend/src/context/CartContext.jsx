import React, { createContext, useState, useContext, useEffect } from 'react';
import axios from 'axios';

const CartContext = createContext();

export const useCart = () => useContext(CartContext);

export const CartProvider = ({ children, cartId }) => {
  const [cartItems, setCartItems] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (cartId) {
      axios.get(`http://localhost:8080/carts/${cartId}`)
        .then(res => {
          if (res.data && res.data.items) {
            setCartItems(res.data.items);
          }
          setLoading(false);
        })
        .catch(err => {
          console.error("Error fetching cart:", err);
          setLoading(false);
        });
    }
  }, [cartId]);

  const addToCart = async (product) => {
    try {
      const response = await axios.post(`http://localhost:8080/carts/${cartId}/items`, {
        product_id: product.ID,
        quantity: 1
      });

      const newItem = response.data;

      setCartItems(prevItems => [...prevItems, newItem]);
    } catch (error) {
      console.error("Error adding to cart:", error);
    }
  };

  const removeFromCart = async (itemId) => {
    try {
      await axios.delete(`http://localhost:8080/carts/${cartId}/items/${itemId}`);
      setCartItems(prevItems => prevItems.filter(item => item.ID !== itemId));
    } catch (error) {
      console.error("Error removing from cart:", error);
    }
  };

  const clearCart = () => {
    setCartItems([]);
  };

  const totalAmount = cartItems.reduce((sum, item) => sum + (item.product.price * item.quantity), 0);

  return (
    <CartContext.Provider value={{ cartItems, addToCart, removeFromCart, clearCart, totalAmount, loading }}>
      {children}
    </CartContext.Provider>
  );
};
