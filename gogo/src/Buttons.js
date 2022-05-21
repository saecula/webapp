import React from "react";
import "./webapp.css";

export const PassButton = ({ onPass, disabled }) => (
  <div className="pass-button">
    <button onClick={onPass} disabled={disabled}>
      pass
    </button>
  </div>
);

export const ResignButton = ({ onResign }) => (
  <div className="resign-button">
    <button onClick={onResign}>resign</button>
  </div>
);

export const SwitchButton = ({ onSwitch }) => (
  <div className="switch-button">
    <button onClick={onSwitch}>swap color</button>
  </div>
);

export default PassButton;
