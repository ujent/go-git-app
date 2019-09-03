import { createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';
import logger from 'redux-logger';

import { rootReducer } from './reducers';

export default function createAppStore(currentUser, currentRepo, currentBranch, repos, branches, files) {
  const initialState = getInitialState(currentUser, currentRepo, currentBranch, repos, branches, files);
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

function getInitialState(currentUser, currentRepo, currentBranch, repos, branches, files) {
  const user = currentUser ? currentUser : '';
  const repo = currentRepo ? currentRepo : '';
  const branch = currentBranch ? currentBranch : '';
  const brItems = branches ? branches : [];
  const repoItems = repos ? repos : [];
  const fileItems = files ? files : [];

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
    currentFile: null,
    // currentFile: {
    //   name: '',
    //   content: '',
    //   isConflict: false
    // },
    files: fileItems,
    confirmPopup: {
      isOpen: false,
      message: '',
      onConfirm: function () { },
      onClose: function () { }
    }

  };

  return initialState;
}
