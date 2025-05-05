export const formatNumber = (value: number): string => {
  if (value === undefined || value === null) {
    return '-';
  }

  // Для больших чисел используем сокращенный формат
  if (value >= 1000000) {
    return `${(value / 1000000).toFixed(2)}M`;
  }
  if (value >= 1000) {
    return `${(value / 1000).toFixed(2)}K`;
  }

  // Для маленьких чисел используем больше знаков после запятой
  if (value < 0.01) {
    return value.toFixed(8);
  }
  if (value < 1) {
    return value.toFixed(6);
  }
  if (value < 10) {
    return value.toFixed(4);
  }

  // Для обычных чисел используем 2 знака после запятой
  return value.toFixed(2);
};

export const formatPrice = (value: number, currency: string = 'USD'): string => {
  return `${formatNumber(value)} ${currency}`;
};

export const formatPercentage = (value: number): string => {
  return `${value >= 0 ? '+' : ''}${value.toFixed(2)}%`;
};

export const formatDate = (date: Date): string => {
  return date.toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}; 