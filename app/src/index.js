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
const store = createAppStore(currentUser);

ReactDOM.render(
    <Provider store={store}><App /></Provider>,
    document.getElementById('root')
);
