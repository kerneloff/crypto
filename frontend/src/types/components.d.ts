export interface OrderBookEntry {
  price: number;
  quantity: number;
}

export interface OrderBookProps {
  symbol: string;
  wsUrl: string;
}

export interface TradingChartProps {
  symbol: string;
  wsUrl: string;
  height?: number;
}

export interface OrderFormProps {
  symbol: string;
  lastPrice: number;
  onSubmit: (order: OrderData) => void;
}

export interface OrderData {
  symbol: string;
  type: 'limit' | 'market';
  side: 'buy' | 'sell';
  price?: number;
  quantity: number;
} 