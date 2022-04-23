import React from "react";
import "./webapp.css";

const ErrorModal = ({ err }) => {
  return (
    <div className="modal-container">
      <div className="error-modal">
        <div
          style={{
            display: "flex",
            flexDirection: "column",
            padding: "10px",
            textAlign: "center",
            margin: "auto",
          }}
        >
          <div style={{ padding: "10px" }}> Cant find that game :(</div>
          <div style={{ padding: "10px" }}> Error: {err.message}</div>
        </div>
      </div>
    </div>
  );
};

export default ErrorModal;
