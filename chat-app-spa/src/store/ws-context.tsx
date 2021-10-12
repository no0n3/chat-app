import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { w3cwebsocket } from "websocket";
import { AuthContext } from "./auth-context";
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

export const WsContext = createContext({
  ws: null,
} as any);

export default function WsContextProvider(props: any) {
  const [socket, setSocket] = useState<any>(null);

  const { token } = useContext(AuthContext);
  const subject = useMemo(() => {
    return new Subject();
  }, []);

  const sendMessage = (type: string, msg: any) => {
    socket.send(buildMessage(type, msg));
  };

  useEffect(() => {
    if (!token) {
      return;
    }

    setSocket(new w3cwebsocket(`${process.env.REACT_APP_WS_ENDPOINT}/ws?token=${token}`))
  }, [token]);

  useEffect(() => {
    if (!socket) {
      return;
    }

    socket.onmessage = (message: any) => {
      const msg: { Type: string, Payload: any } = strToJson(message.data as string);

      subject.next(msg);
    };
  }, [socket]);

  return (
    <WsContext.Provider value={{
      sendMessage,
      subject
    }}>
      {props.children}
    </WsContext.Provider>
  );
};
