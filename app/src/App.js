import React from 'react';
import './App.css';
import Settings from './containers/SettingsContainer'
import Commands from './containers/CommandsContainer'
import Output from './components/Output'


const App = props => {
  return (
    <div className="app-wrapper">
      <header className="App-header">
      </header>
      <div className="app-main">
        <Commands />
        <div className="main-content">
          <Settings />
          <Output message={props.outputMsg} />
          <section className="repo-content">
            <h2 className="visually-hidden">Repository content</h2>
            <div className="repo-files">
              <h3>Files</h3>
              <ul className="files-list">
                <li>file1</li>
              </ul>
            </div>
            <div className="file-content ">
              <h3>File content</h3>
              <div className="file-content-buttons">
                <button type="button" className="button">Save</button>
              </div>
              <textarea rows="32"></textarea>
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}

export default App;
