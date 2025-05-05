import React, { useEffect, useState } from 'react';
import { Box, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Typography } from '@mui/material';
import { formatNumber } from '../utils/format';
import { OrderBookProps, OrderBookEntry } from '../types/components';

export const OrderBook: React.FC<OrderBookProps> = ({ symbol, wsUrl }) => {
  const [bids, setBids] = useState<OrderBookEntry[]>([]);
  const [asks, setAsks] = useState<OrderBookEntry[]>([]);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    const socket = new WebSocket(wsUrl);
    setWs(socket);

    socket.onopen = () => {
      socket.send(JSON.stringify({
        type: 'subscribe',
        channel: `orderbook.${symbol}`
      }));
    };

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === 'message' && data.channel === `orderbook.${symbol}`) {
        const orderBook = data.data;
        setBids(orderBook.bids);
        setAsks(orderBook.asks);
      }
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    return () => {
      if (socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({
          type: 'unsubscribe',
          channel: `orderbook.${symbol}`
        }));
        socket.close();
      }
    };
  }, [symbol, wsUrl]);

  const renderOrderBookSide = (orders: OrderBookEntry[], isBids: boolean) => {
    return orders.map((order, index) => (
      <TableRow key={index}>
        <TableCell sx={{ color: isBids ? 'success.main' : 'error.main' }}>
          {formatNumber(order.price)}
        </TableCell>
        <TableCell>{formatNumber(order.quantity)}</TableCell>
        <TableCell>{formatNumber(order.price * order.quantity)}</TableCell>
      </TableRow>
    ));
  };

  return (
    <Box sx={{ width: '100%', maxWidth: 600, mx: 'auto', my: 2 }}>
      <Typography variant="h6" gutterBottom>
        Order Book - {symbol}
      </Typography>
      <TableContainer component={Paper}>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Price</TableCell>
              <TableCell>Quantity</TableCell>
              <TableCell>Total</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {renderOrderBookSide(asks.slice().reverse(), false)}
            <TableRow>
              <TableCell colSpan={3} align="center" sx={{ bgcolor: 'grey.100' }}>
                Spread: {asks[0] && bids[0] ? formatNumber(asks[0].price - bids[0].price) : '-'}
              </TableCell>
            </TableRow>
            {renderOrderBookSide(bids, true)}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  );
}; 