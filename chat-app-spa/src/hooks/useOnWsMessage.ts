import { useContext, useEffect } from "react";
import { MessageType, WsContext } from "../store/WsContext";

export default function useOnWsMessage(
  callback: (msg: MessageType) => void,
  dependencies: any[] = []
) {
  const { subject } = useContext(WsContext);

  useEffect(() => {
    const subscription = subject.subscribe(callback);

    return () => {
      subscription.unsubscribe();
    };
  }, [...dependencies]);
}
