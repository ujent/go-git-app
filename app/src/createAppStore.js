import { createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';
import logger from 'redux-logger';

import { rootReducer } from './reducers';

export default function createAppStore(currentUser, currentRepo, currentBranch, repos, branches) {
  const initialState = getInitialState(currentUser, currentRepo, currentBranch, repos, branches);
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

function getInitialState(currentUser, currentRepo, currentBranch, repos, branches) {

  const user = currentUser ? currentUser : '';
  const repo = currentRepo ? currentRepo : '';
  const branch = currentBranch ? currentBranch : '';
  const brItems = branches ? branches : [];
  const repoItems = repos ? repos : [];

  const initialState = {
    users: ['user1', 'user2', 'user3'],
    repositories: repoItems,
    branches: brItems,
    settings: {
      currentUser: user,
      currentRepo: repo,
      currentBranch: branch
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

function getInitialStateTest() {

  const initialState = {
    users: ['user1', 'user2', 'user3'],
    repositories: ['repo1', 'repo2'],
    branches: ['branch1', 'branch2'],
    settings: {
      currentUser: 'user1',
      currentRepo: 'repo1',
      currentBranch: 'branch1'
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
