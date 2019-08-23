import * as api from './api';
import { ActionType } from './constants';

export function showMessage(msg) {
    return {
        type: ActionType.SHOW_MSG,
        msg
    };
}

export function resetMessage() {
    return {
        type: ActionType.RESET_MSG,
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

export function setUsers(users, current) {
    return {
        type: ActionType.SET_USERS,
        users,
        current
    };
}

export function getUsers() {
    return (dispatch, getState) => {
        api.getUsers().then(
            rs => {
                dispatch(setUsers(rs.users, rs.current));
            },
            err => {
                dispatch(showError(err));
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
                dispatch(showError(err));
            }
        );
    }
}

export function setRepositories(repos, current) {
    return {
        type: ActionType.SET_REPOSITORIES,
        repos,
        current
    };
}

export function getRepositories() {
    return (dispatch, getState) => {
        api.getRepositories().then(
            rs => {
                dispatch(setRepositories(rs.repos, rs.current));
            },
            err => {
                dispatch(showError(err));
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
                dispatch(showError(err));
            }
        );
    }
}

export function setBranches(branches, current) {
    return {
        type: ActionType.SET_BRANCHES,
        branches,
        current
    };
}

export function getBranches() {
    return (dispatch, getState) => {
        api.getBranches().then(
            rs => {
                dispatch(setBranches(rs.branches, rs.current));
                dispatch(getRepoFiles())
            },
            err => {
                dispatch(showError(err));
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
                dispatch(showError(err));
            }
        );
    }
}

export function getRepoFiles() {
    return (dispatch, getState) => {
        api.getRepoFiles().then(
            files => dispatch(setFiles(files)),
            err => {
                dispatch(showError(err));
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