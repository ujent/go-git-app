import React from 'react';
import './App.css';
import Settings from './containers/SettingsContainer';
import Commands from './containers/CommandsContainer';
import RepoContent from './containers/RepoContentContainer';
import Output from './components/Output';
import Spinner from './components/Spinner';
import Popup from './containers/ConfirmPopupContainer';


const App = props => {
  return (
    <>
      <div className="app-wrapper">
        <header className="App-header">
        </header>
        <div className="app-main">
          <Commands />
          <div className="main-content">
            <Settings />
            <Output message={props.outputMsg} />
            <RepoContent />
          </div>
        </div>
      </div>
      <Popup />
      <Spinner isVisible={props.isSpinnerVisible} />
    </>
  );
}

export default App;
