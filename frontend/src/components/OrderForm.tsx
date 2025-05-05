import React, { useState } from 'react';
import {
  Box,
  Button,
  FormControl,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  TextField,
  Typography,
  SelectChangeEvent,
} from '@mui/material';
import { formatNumber } from '../utils/format';
import { OrderFormProps, OrderData } from '../types/components';

export const OrderForm: React.FC<OrderFormProps> = ({
  symbol,
  lastPrice,
  onSubmit,
}) => {
  const [orderType, setOrderType] = useState<'limit' | 'market'>('limit');
  const [side, setSide] = useState<'buy' | 'sell'>('buy');
  const [price, setPrice] = useState<string>(lastPrice.toString());
  const [quantity, setQuantity] = useState<string>('');

  const handleTypeChange = (event: SelectChangeEvent) => {
    setOrderType(event.target.value as 'limit' | 'market');
  };

  const handleSideChange = (event: SelectChangeEvent) => {
    setSide(event.target.value as 'buy' | 'sell');
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const orderData: OrderData = {
      symbol,
      type: orderType,
      side,
      quantity: parseFloat(quantity),
    };

    if (orderType === 'limit') {
      orderData.price = parseFloat(price);
    }

    onSubmit(orderData);
  };

  const handlePriceChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    if (/^\d*\.?\d*$/.test(value)) {
      setPrice(value);
    }
  };

  const handleQuantityChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    if (/^\d*\.?\d*$/.test(value)) {
      setQuantity(value);
    }
  };

  return (
    <Box sx={{ width: '100%', maxWidth: 400, mx: 'auto', my: 2 }}>
      <Paper sx={{ p: 2 }}>
        <Typography variant="h6" gutterBottom>
          New Order
        </Typography>
        <form onSubmit={handleSubmit}>
          <FormControl fullWidth margin="normal">
            <InputLabel>Order Type</InputLabel>
            <Select
              value={orderType}
              label="Order Type"
              onChange={handleTypeChange}
            >
              <MenuItem value="limit">Limit</MenuItem>
              <MenuItem value="market">Market</MenuItem>
            </Select>
          </FormControl>

          <FormControl fullWidth margin="normal">
            <InputLabel>Side</InputLabel>
            <Select
              value={side}
              label="Side"
              onChange={handleSideChange}
            >
              <MenuItem value="buy">Buy</MenuItem>
              <MenuItem value="sell">Sell</MenuItem>
            </Select>
          </FormControl>

          {orderType === 'limit' && (
            <TextField
              fullWidth
              margin="normal"
              label="Price"
              value={price}
              onChange={handlePriceChange}
              type="text"
              inputProps={{
                inputMode: 'decimal',
                pattern: '^\\d*\\.?\\d*$',
              }}
            />
          )}

          <TextField
            fullWidth
            margin="normal"
            label="Quantity"
            value={quantity}
            onChange={handleQuantityChange}
            type="text"
            inputProps={{
              inputMode: 'decimal',
              pattern: '^\\d*\\.?\\d*$',
            }}
          />

          <Box sx={{ mt: 2, display: 'flex', justifyContent: 'space-between' }}>
            <Typography variant="body2">
              Last Price: {formatNumber(lastPrice)}
            </Typography>
            {orderType === 'limit' && price && quantity && (
              <Typography variant="body2">
                Total: {formatNumber(parseFloat(price) * parseFloat(quantity))}
              </Typography>
            )}
          </Box>

          <Button
            fullWidth
            variant="contained"
            color={side === 'buy' ? 'success' : 'error'}
            type="submit"
            sx={{ mt: 2 }}
          >
            {side === 'buy' ? 'Buy' : 'Sell'} {symbol}
          </Button>
        </form>
      </Paper>
    </Box>
  );
}; 