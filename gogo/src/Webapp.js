import React, { useState, useEffect, useCallback } from "react";
import Board from "./Board";
import "./webapp.css";

const Webapp = () => {
  const [socket, setSocket] = useState(null);

  useEffect(() => {
    initSocket();
    return disconnectSocket;
  }, []);

  const initSocket = useCallback(() => {
    const s = new WebSocket(`ws://localhost:4000/ws`);
    setSocket(s);
  }, [socket]);

  const disconnectSocket = () => {
    socket?.disconnect && socket.disconnect();
    setSocket(null);
  };

  return (
    <div className="webapp">
      <header className="webapp-header">
        <Board socket={socket} />
      </header>
    </div>
  );
};

export default Webapp;
