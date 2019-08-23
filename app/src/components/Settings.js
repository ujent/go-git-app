/*import React from 'react';
import '../App.css';
import './LoginForm.css';

const LoginForm = props => {
  function handleChangeLogin(e) {
    props.handleChangeLoginForm(e.target.value, props.password);
  }

  function handleChangePsw(e) {
    props.handleChangeLoginForm(props.login, e.target.value);
  }

  function handleEnterPress(e) {
    if (e.key === 'Enter') {
      props.onSubmit(props.login, props.password);
    }
  }

  return (
    <div
      className="App login popup-overlay"
      onKeyPress={e => handleEnterPress(e)}
    >
      <div className="popup" onClick={e => e.stopPropagation()}>
        <input
          type="text"
          placeholder="Login"
          value={props.login}
          onChange={e => handleChangeLogin(e)}
        />
        <input
          type="password"
          placeholder="Password"
          value={props.password}
          onChange={e => handleChangePsw(e)}
        />
        <div className="validate-error">{props.message}</div>
        <div>
          <button
            disabled={!props.login || !props.password}
            className="popup-close-btn"
            onClick={e => props.onSubmit(props.login, props.password)}
          >
            Submit
          </button>
        </div>
      </div>
    </div>
  );
};

export default LoginForm;*/
