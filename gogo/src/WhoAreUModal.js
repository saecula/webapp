import React, { useCallback } from "react";
import "./webapp.css";

const WhoAreUModal = () => {
  const handleSubmit = useCallback((e) => {
    e.preventDefault();
    console.log("yey", e.target[0]?.value);
  }, []);
  return (
    <div className="modal-container">
      <div className="modal">
        who are u?
        <form onSubmit={handleSubmit}>
          <input type="text" className="modal-input" />
        </form>
      </div>
    </div>
  );
};

export default WhoAreUModal;
