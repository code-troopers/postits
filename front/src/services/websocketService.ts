import { useBoardStore } from '@/stores/board';
import keycloak from '../keycloak';
import type { Message } from '../models/Message';

let ws: WebSocket | null = null;
const WS_URL = import.meta.env.VITE_WS_URL;

export function connectWebSocket(): void {
  const token = keycloak.token;

  if (!token) {
    console.error('No Keycloak token found');
    return;
  }

  // Inclure le token dans la chaîne de requête
  const wsUrl = `${WS_URL}/ws?token=${encodeURIComponent(token)}`;
  ws = new WebSocket(wsUrl);

  ws.onopen = () => {
    console.log('WebSocket connection established');
  };

  ws.onmessage = (event: MessageEvent) => {
    const store = useBoardStore();
    const message: Message = JSON.parse(event.data);
    store.onMessage(message);
    console.log('Message received:', message);
  };

  ws.onerror = (error: Event) => {
    console.error('WebSocket error:', error);
  };

  ws.onclose = () => {
    console.log('WebSocket connection closed');
  };
}

export function sendMessage(message: Message): void {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(message));
  } else {
    throw new Error('WebSocket is not connected');
  }
}
