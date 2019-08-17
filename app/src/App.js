import React from 'react';
import './App.css';

function App() {
  return (
    <div className="app-wrapper">
      <header className="App-header">
      </header>
      <body>
        <div className="app-main">
          <section className="commands">
            <h2>Commands</h2>
            <ul className="commands-list">
              <li>
                <button type="button" className="button">Commit</button>
                <input type="text" placeholder="message" />
                <hr></hr>
              </li>
              <li>
                <button type="button" className="button">Merge</button>
                <select placeholder="branch">
                  <option value="branch1">branch1</option>
                  <option value="branch2">branch2</option>
                </select>
                <hr></hr>
              </li>
              <li>
                <div>
                  <button type="button" className="button">Clone</button>
                  <input type="text" placeholder="URL" />
                </div>
                <div>
                  <input type="text" placeholder="user" />
                  <input type="password" placeholder="password" />
                </div>
                <hr></hr>
              </li>
              <li>
                <button type="button" className="button">Remove repo</button>
                <input type="text" placeholder="name" />
                <hr></hr>
              </li>
              <li>
                <div>
                  <button type="button" className="button">Pull</button>
                  <input type="text" placeholder="remote" />
                </div>
                <div>
                  <input type="text" placeholder="user" />
                  <input type="password" placeholder="password" />
                </div>
                <hr></hr>
              </li>
              <li>
                <div>
                  <button type="button" className="button">Push</button>
                  <input type="text" placeholder="remote" />
                </div>
                <div>
                  <input type="text" placeholder="user" />
                  <input type="password" placeholder="password" />
                </div>
                <hr></hr>
              </li>
              <li>
                <button type="button" className="button">Checkout branch</button>
                <select placeholder="branch">
                  <option value="branch1">branch1</option>
                  <option value="branch2">branch2</option>
                </select>
                <hr></hr>
              </li>
              <li>
                <button type="button" className="button">Remove branch</button>
                <select placeholder="branch">
                  <option value="branch1">branch1</option>
                  <option value="branch2">branch2</option>
                </select>
                <hr></hr>
              </li>
              <li>
                <button type="button" className="button">Log</button>
              </li>
            </ul>
          </section>
          <div className="main-content">
            <section className="main-settings">
              <label htmlFor="userSelectID">User</label>
              <select id="userSelectID">
                <option value="user1">user1</option>
                <option value="user2">user2</option>
              </select>
              <label htmlFor="repoSelectID">Repository</label>
              <select id="repoSelectID">
                <option value="repo1">repo1</option>
                <option value="repo2">repo2</option>
              </select>
              <label htmlFor="branchSelectID">Branch</label>
              <select id="branchSelectID">
                <option value="branch1">branch1</option>
                <option value="branch2">branch2</option>
              </select>
            </section>
            <section className="output">
              <p></p>
            </section>
            <section className="repo-content">
              <h2 className="visually-hidden">Repository content</h2>
              <div className="files-list">
                <h3>Files</h3>
                <ul>
                  <li>file1</li>
                </ul>
              </div>
              <div className="file-content">
                <h3>File content</h3>
                <div>
                  <button type="button" className="button">Reset</button>
                  <button type="button" className="button">Save</button>
                </div>
                <textarea cols="30" rows="10"></textarea>
              </div>
            </section>
          </div>
        </div>


      </body>
    </div>
  );
}

export default App;
