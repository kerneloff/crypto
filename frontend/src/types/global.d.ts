declare module 'lightweight-charts' {
  export interface IChartApi {
    applyOptions(options: any): void;
    remove(): void;
    addCandlestickSeries(options?: any): any;
  }

  export enum ColorType {
    Solid = 'solid',
  }

  export function createChart(container: HTMLElement, options?: any): IChartApi;
}

declare namespace JSX {
  interface IntrinsicElements {
    [elemName: string]: any;
  }
} 