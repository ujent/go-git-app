import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import './index.css';
import App from './containers/AppContainer';
import * as serviceWorker from './serviceWorker';
import createAppStore from './createAppStore';


// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();


/* ReactDOM.render(<App />, document.getElementById('root')); */
const currentUser = window.localStorage.getItem('go_git_user');
const currentRepo = window.localStorage.getItem('go_git_repo');
const currentBranch = window.localStorage.getItem('go_git_branch');
const store = createAppStore(currentUser, currentRepo, currentBranch);

ReactDOM.render(
    <Provider store={store}><App /></Provider>,
    document.getElementById('root')
);
