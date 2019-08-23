import React from 'react';
import './App.css';
import Output from './components/Output'

function App() {
  return (
    <div className="app-wrapper">
      <header className="App-header">
      </header>
      <div className="app-main">
        <section className="commands">
          <h2>Commands</h2>
          <ul className="commands-list">
            <li className="commit-command">
              <button type="button" className="button">Commit</button>
              <textarea placeholder="commit message" rows="3"></textarea>
              <hr></hr>
            </li>
            <li className="command-block">
              <div className="command-block-content">
                <button type="button" className="button">Merge</button>
                <select placeholder="branch" className="command-block-input" defaultValue="">
                  <option value="" disabled hidden>select</option>
                  <option value="branch1">branch1</option>
                  <option value="branch2">branch2</option>
                </select></div>

              <hr></hr>
            </li>
            <li>
              <div className="command-block command-block-content">
                <button type="button" className="button">Clone</button>
                <input type="text" placeholder="URL" className="command-block-input" />
              </div>
              <div className="credentials">
                <input type="text" placeholder="user" />
                <input type="password" placeholder="password" />
              </div>
              <hr></hr>
            </li>
            <li>
              <div className="command-block command-block-content">
                <button type="button" className="button">Pull</button>
                <input type="text" placeholder="remote" className="command-block-input" />
              </div>
              <div className="credentials">
                <input type="text" placeholder="user" />
                <input type="password" placeholder="password" />
              </div>
              <hr></hr>
            </li>
            <li>
              <div className="command-block command-block-content">
                <button type="button" className="button">Push</button>
                <input type="text" placeholder="remote" className="command-block-input" />
              </div>
              <div className="credentials">
                <input type="text" placeholder="user" />
                <input type="password" placeholder="password" />
              </div>
              <hr></hr>
            </li>
            <li className="command-block">
              <div className="command-block-content">
                <button type="button" className="button medium-button">Remove repo</button>
                <select placeholder="select" className="command-block-input" defaultValue="">
                  <option value="" disabled hidden>select</option>
                  <option value="repo1">repo1</option>
                  <option value="repo2">repo2</option>
                </select>
              </div>
              <hr></hr>
            </li>
            <li className="command-block">
              <div className="command-block-content">
                <button type="button" className="button medium-button">Checkout branch</button>
                <select placeholder="select" className="command-block-input" defaultValue="">
                  <option value="" disabled hidden>select</option>
                  <option value="branch1">branch1</option>
                  <option value="branch2">branch2</option>
                </select>
              </div>
              <hr></hr>
            </li>
            <li className="command-block">
              <div className="command-block-content">
                <button type="button" className="button medium-button">Remove branch</button>
                <select placeholder="select" className="command-block-input" defaultValue="">
                  <option value="" disabled hidden>select</option>
                  <option value="branch1">branch1</option>
                  <option value="branch2">branch2</option>
                </select>
              </div>
              <hr></hr>
            </li>
            <li>
              <button type="button" className="button">Log</button>
            </li>
          </ul>
        </section>
        <div className="main-content">
          <section className="main-settings">
            <div>
              <label htmlFor="userSelectID">User</label>
              <select id="userSelectID" defaultValue="">
                <option value="" disabled hidden>select</option>
                <option value="user1">user1</option>
                <option value="user2">user2</option>
                <option value="user3">user3</option>
              </select>
            </div>
            <div>
              <label htmlFor="repoSelectID">Repository</label>
              <select id="repoSelectID" defaultValue="">
                <option value="" disabled hidden>select</option>
                <option value="repo1">repo1</option>
                <option value="repo2">repo2</option>
              </select>
            </div>
            <div>
              <label htmlFor="branchSelectID">Branch</label>
              <select id="branchSelectID" defaultValue="">
                <option value="" disabled hidden>select</option>
                <option value="branch1">branch1</option>
                <option value="branch2">branch2</option>
              </select>
            </div>
          </section>
          <Output message="" />
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
