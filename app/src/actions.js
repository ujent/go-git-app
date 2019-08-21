import * as api from './api';
import { ActionType } from './constants';

export function showMessage(msg) {
    return {
        type: ActionType.SHOW_MSG,
        msg
    };
}

export function showError(err) {
    const msg = `Error code: ${err.status}
    Message: ${err.message}`

    return {
        type: ActionType.SHOW_MSG,
        msg
    };
}

export function setUsers(users) {
    return {
        type: ActionType.SET_USERS,
        users
    };
}

export function getUsers() {
    return (dispatch, getState) => {
        api.getUsers().then(
            users => {
                dispatch(setUsers(users));
            },
            err => {
                dispatch(showMessage(err.message));
            }
        );
    }
}

export function switchUser(name) {
    return (dispatch, getState) => {
        api.switchUser(name).then(
            () => {
                dispatch(resetRepo)
                dispatch(resetBranch)
            },
            err => {
                dispatch(showMessage(err.message));
            }
        );
    }
}

export function setRepositories(repos) {
    return {
        type: ActionType.SET_REPOSITORIES,
        repos
    };
}

export function getRepositories() {
    return (dispatch, getState) => {
        api.getRepositories().then(
            repos => {
                dispatch(setRepositories(repos));
            },
            err => {
                dispatch(showMessage(err.message));
            }
        );
    }
}

export function resetRepo() {
    return {
        type: ActionType.RESET_REPO
    };
}

export function switchRepo(name) {
    return (dispatch, getState) => {
        api.switchRepo(name).then(
            () => { dispatch(getBranches) },
            err => {
                dispatch(showMessage(err.message));
            }
        );
    }
}

export function setBranches(branches) {
    return {
        type: ActionType.SET_BRANCHES,
        branches
    };
}

export function getBranches() {
    return (dispatch, getState) => {
        api.getBranches().then(
            branches => {
                dispatch(setBranches(branches));
                dispatch(getRepoFiles())
            },
            err => {
                dispatch(showMessage(err.message));
            }
        );
    }
}

export function resetBranch() {
    return {
        type: ActionType.RESET_BRANCH
    };
}

export function switchBranch(name) {
    return (dispatch, getState) => {
        api.switchBranch(name).then(
            () => dispatch(getRepoFiles()),
            err => {
                dispatch(showMessage(err.message));
            }
        );
    }
}

export function getRepoFiles() {
    return (dispatch, getState) => {
        api.getRepoFiles().then(
            files => dispatch(setFiles(files)),
            err => {
                dispatch(showMessage(err.message));
            }
        );
    }
}

export function setFiles(files) {
    return {
        type: ActionType.SET_FILES,
        files
    };
}

/* export function changeBenefitsFilter(newFilter) {
  return {
    type: ActionType.CHANGE_BENEFITS_FILTER,
    newFilter
  };
}

export function startFreeco() {
  return (dispatch, getState) => {
    const user = getUser(getState);
    api.startFreeco(user).then(
      () => {
        dispatch(getProcessInfo());
      },
      err => {
        dispatch(showError(err));
      }
    );
  };
}*/