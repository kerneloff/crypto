import React, { useState } from 'react';
import { Container, Grid, Paper, Typography } from '@mui/material';
import { OrderBook } from '../components/OrderBook';
import { TradingChart } from '../components/TradingChart';
import { OrderForm } from '../components/OrderForm';
import { OrderData } from '../types/components';

const Home: React.FC = () => {
  const [selectedSymbol] = useState('BTC/USDT');
  const [lastPrice, setLastPrice] = useState(50000);

  const handleOrderSubmit = async (orderData: OrderData) => {
    try {
      const response = await fetch('/api/v1/private/orders', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(orderData),
      });

      if (!response.ok) {
        throw new Error('Failed to create order');
      }

      // Обновляем последнюю цену
      const data = await response.json();
      if (data.price) {
        setLastPrice(data.price);
      }
    } catch (error) {
      console.error('Error creating order:', error);
    }
  };

  return (
    <Container maxWidth="xl">
      <Typography variant="h4" component="h1" gutterBottom sx={{ mt: 4 }}>
        Crypto Exchange
      </Typography>

      <Grid container spacing={3}>
        <Grid item xs={12} md={8}>
          <Paper sx={{ p: 2 }}>
            <TradingChart
              symbol={selectedSymbol}
              wsUrl="ws://localhost:8080/ws"
              height={500}
            />
          </Paper>
        </Grid>

        <Grid item xs={12} md={4}>
          <OrderForm
            symbol={selectedSymbol}
            lastPrice={lastPrice}
            onSubmit={handleOrderSubmit}
          />
        </Grid>

        <Grid item xs={12}>
          <Paper sx={{ p: 2 }}>
            <OrderBook
              symbol={selectedSymbol}
              wsUrl="ws://localhost:8080/ws"
            />
          </Paper>
        </Grid>
      </Grid>
    </Container>
  );
};

export default Home; 