import { useContext, useEffect, useState } from "react";
import { w3cwebsocket } from "websocket";
import { AuthContext } from "../store/auth-context";

export function useWebsocket() {
  const [socket, setSocket] = useState<any>(null);

  const { token } = useContext(AuthContext)

  useEffect(() => {
    if (!token) {
      if (socket) {
        socket.close();
      }

      return;
    }

    if (!socket) {
      setSocket(new w3cwebsocket(`${process.env.REACT_APP_WS_ENDPOINT}/ws?token=${token}`))
    }
  }, [token, socket]);

  return socket;
}
