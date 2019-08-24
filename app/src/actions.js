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

export function checkoutBranch() {
    return (dispatch, getState) => {
        api.checkoutBranch().then(
            rs => {

            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}

export function clone() {
    return (dispatch, getState) => {
        api.clone().then(
            rs => {

            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}
export function commit() {
    return (dispatch, getState) => {
        api.commit().then(
            rs => {

            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}
export function log() {
    return (dispatch, getState) => {
        api.log().then(
            rs => {

            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}
export function merge() {
    return (dispatch, getState) => {
        api.merge().then(
            rs => {

            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}
export function pull() {
    return (dispatch, getState) => {
        api.pull().then(
            rs => {

            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}
export function push() {
    return (dispatch, getState) => {
        api.push().then(
            rs => {

            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}
export function removeBranch() {
    return (dispatch, getState) => {
        api.removeBranch().then(
            rs => {

            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}
export function removeRepo() {
    return (dispatch, getState) => {
        api.removeRepo().then(
            rs => {

            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}