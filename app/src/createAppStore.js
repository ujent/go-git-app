import { createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';
import logger from 'redux-logger';

import { rootReducer } from './reducers';

export default function createAppStore() {
  const initialState = getInitialState();
  let store;

  if (process.env.NODE_ENV === 'production') {
    store = createStore(rootReducer, initialState, applyMiddleware(thunk));
  } else {
    store = createStore(
      rootReducer,
      initialState,
      applyMiddleware(thunk, logger)
    );
  }

  return store;
}

function getInitialState() {


  const initialState = {
    users: ['user1', 'user2', 'user3'],
    repositories: [],
    branches: [],
    settings: {
      currentUser: 'user1',
      currentRepo: 'repo1',
      currentBranch: ''
    },
    output: '',
    fileContent: {
      isVisible: false,
      content: ''
    },
    files: []

  };

  return initialState;
}
