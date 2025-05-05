import React, { useEffect, useRef } from 'react';
import { Box, Paper, Typography } from '@mui/material';
import { createChart, ColorType, IChartApi } from 'lightweight-charts';

interface TradingChartProps {
  symbol: string;
  wsUrl: string;
  height?: number;
}

interface CandleData {
  time: string;
  open: number;
  high: number;
  low: number;
  close: number;
}

export const TradingChart: React.FC<TradingChartProps> = ({
  symbol,
  wsUrl,
  height = 400,
}) => {
  const chartContainerRef = useRef<HTMLDivElement>(null);
  const chartRef = useRef<IChartApi | null>(null);
  const candlestickSeriesRef = useRef<any>(null);

  useEffect(() => {
    if (chartContainerRef.current) {
      const chart = createChart(chartContainerRef.current, {
        layout: {
          background: { type: ColorType.Solid, color: '#ffffff' },
          textColor: '#333',
        },
        grid: {
          vertLines: { color: '#f0f0f0' },
          horzLines: { color: '#f0f0f0' },
        },
        width: chartContainerRef.current.clientWidth,
        height,
      });

      const candlestickSeries = chart.addCandlestickSeries({
        upColor: '#26a69a',
        downColor: '#ef5350',
        borderVisible: false,
        wickUpColor: '#26a69a',
        wickDownColor: '#ef5350',
      });

      chartRef.current = chart;
      candlestickSeriesRef.current = candlestickSeries;

      const handleResize = () => {
        if (chartContainerRef.current) {
          chart.applyOptions({
            width: chartContainerRef.current.clientWidth,
          });
        }
      };

      window.addEventListener('resize', handleResize);

      return () => {
        window.removeEventListener('resize', handleResize);
        chart.remove();
      };
    }
  }, [height]);

  useEffect(() => {
    const socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      socket.send(JSON.stringify({
        type: 'subscribe',
        channel: `trades.${symbol}`,
      }));
    };

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === 'message' && data.channel === `trades.${symbol}`) {
        const trade = data.data;
        if (candlestickSeriesRef.current) {
          // Здесь должна быть логика агрегации сделок в свечи
          // Для примера просто добавляем точку
          candlestickSeriesRef.current.update({
            time: new Date(trade.timestamp).toISOString(),
            open: trade.price,
            high: trade.price,
            low: trade.price,
            close: trade.price,
          });
        }
      }
    };

    return () => {
      if (socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({
          type: 'unsubscribe',
          channel: `trades.${symbol}`,
        }));
        socket.close();
      }
    };
  }, [symbol, wsUrl]);

  return (
    <Box sx={{ width: '100%', my: 2 }}>
      <Typography variant="h6" gutterBottom>
        {symbol} Chart
      </Typography>
      <Paper>
        <div ref={chartContainerRef} />
      </Paper>
    </Box>
  );
}; 