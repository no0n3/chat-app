import { createContext, useContext, useEffect, useMemo, useRef, useState } from "react";
import { w3cwebsocket } from "websocket";
import { AuthContext } from "./AuthContext";
import { Subject } from 'rxjs';
import { buildMessage } from '../WsUtils';

const parseJSON = (jsonStr: string) => {
  try {
    return JSON.parse(jsonStr);
  } catch (e) {
    return {}
  }
}

const strToJson = (str: string) => {
  try {
    return parseJSON(str);
  } catch (e) {
    return {}
  }
}

export type MessageType = {
  Type: string;
  Payload: any;
};

export const WsContext = createContext({
  ws: null,
} as any);

export default function WsContextProvider(props: any) {
  const wsRef = useRef<w3cwebsocket>();

  const { token } = useContext(AuthContext);
  const subject = useMemo(() => new Subject<MessageType>(), []);

  const sendMessage = (type: string, msg: MessageType) => {
    if (!wsRef.current) return;

    wsRef.current.send(buildMessage(type, msg));
  };

  const createWs = () => {
    const ws = new w3cwebsocket(`${process.env.REACT_APP_WS_ENDPOINT}/ws?token=${token}`);

    let pongTimeout: NodeJS.Timeout;

    ws.onmessage = (message: any) => {
      const msg: MessageType = strToJson(message.data as string);

      if (!msg?.Type) return;

      if (msg.Type === 'pong') {
        clearTimeout(pongTimeout);
        ping();

        return;
      }

      subject.next(msg);
    };

    ws.onerror = () => {
      setTimeout(() => {
        createWs();
      }, 15000);
    };

    const ping = () => {
      pongTimeout = setTimeout(() => {
        ws.close();

        createWs();
      }, 15000);
    };

    wsRef.current = ws;
  };

  useEffect(() => {
    if (!token) {
      return;
    }

    if (wsRef.current) {
      wsRef.current.close();
    }

    createWs();
  }, [token]);

  return (
    <WsContext.Provider value={{
      sendMessage,
      subject
    }}>
      {props.children}
    </WsContext.Provider>
  );
};
