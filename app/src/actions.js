import * as api from './api';
import { ActionType } from './constants';


export function getSettings(state) {
    return { user: state.settings.currentUser, repo: state.settings.currentRepo, branch: state.settings.currentBranch }
}

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

export function setCurrentUser(name) {
    return {
        type: ActionType.SET_CURRENT_USER,
        user: name
    }
}

export function switchUser(name) {
    return (dispatch, getState) => {
        api.switchUser(name).then(
            () => {
                dispatch(setCurrentUser(name));
                dispatch(getRepositories(name))
            },
            err => {
                dispatch(showError(err));
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


export function setCurrentRepo(repo) {
    return {
        type: ActionType.SET_CURRENT_REPOSITORY,
        current: repo
    };
}

export function getRepositories(user) {
    return (dispatch, getState) => {
        //const user = getState().currentUser;

        api.getRepositories(user).then(
            rs => {
                dispatch(setRepositories(rs.repos));
            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}

export function switchRepo(name) {
    return (dispatch, getState) => {
        const user = getState().settings.currentUser;

        api.switchRepo(user, name).then(
            () => {
                dispatch(setCurrentRepo(name));
                dispatch(getBranches(user, name));
            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}

export function setBranches(branches, currentBranch) {
    return {
        type: ActionType.SET_BRANCHES,
        branches,
        current: currentBranch
    };
}

export function setCurrentBranch(branch) {
    return {
        type: ActionType.SET_CURRENT_BRANCH,
        current: branch
    };
}

export function getBranches(user, repo) {
    return (dispatch, getState) => {
        api.getBranches(user, repo).then(
            rs => {
                dispatch(setBranches(rs.branches, rs.current));

                if (rs.currentBranch) {
                    dispatch(getRepoFiles(user, repo, rs.currentBranch));
                }
            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}

export function switchBranch(name) {
    return (dispatch, getState) => {
        const settings = getSettings(getState())

        api.checkoutBranch(settings.user, settings.repo, name).then(
            () => {
                dispatch(setCurrentBranch(name));
                dispatch(getRepoFiles(settings.user, settings.repo, name));
            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}

export function getRepoFiles(user, repo, branch) {
    return (dispatch, getState) => {

        api.getRepoFiles(user, repo, branch).then(
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

export function removeRepoEntry(repo) {
    return {
        type: ActionType.REMOVE_REPO,
        repo
    };
}

export function removeBranchEntry(branch) {
    return {
        type: ActionType.REMOVE_BRANCH,
        branch
    };
}

export function clone(url, authName, authPsw) {
    return (dispatch, getState) => {
        const user = getState().settings.currentUser;

        api.clone(user, url, authName, authPsw).then(
            (rs) => {
                dispatch(switchRepo(rs.name));
                dispatch(showMessage("Success"));
            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}
export function commit(msg) {
    return (dispatch, getState) => {
        const settings = getSettings(getState())

        api.commit(settings, msg).then(

            () => {
                dispatch(showMessage("Success"));
            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}
export function log() {
    return (dispatch, getState) => {
        const settings = getSettings(getState())

        api.log(settings.user, settings.repo, settings.branch).then(
            rs => {
                let res
                rs.commits.forEach(function (el) {
                    res = res + `commit ${el.hash}
                    Author: ${el.author ? el.author.name : ''} ${el.author ? el.author.email : ''}
                    Date: ${el.date}
                    Msg: ${el.msg}
                        `
                });

                dispatch(showMessage(res))
            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}

export function pull(remote, authName, authPsw) {
    return (dispatch, getState) => {
        const settings = getSettings(getState())

        api.pull(settings, remote, authName, authPsw).then(
            () => {
                dispatch(showMessage("Success"));
            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}
export function push(remote, authName, authPsw) {
    return (dispatch, getState) => {
        const settings = getSettings(getState())

        api.push(settings, remote, authName, authPsw).then(
            () => {
                dispatch(showMessage("Success"));
            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}
export function removeBranch(branch) {
    return (dispatch, getState) => {
        const settings = getSettings(getState())

        api.removeBranch(settings.user, settings.repo, branch).then(
            () => {
                dispatch(showMessage("Success"));
                dispatch(removeBranchEntry(branch));
            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}
export function removeRepo(repo) {
    return (dispatch, getState) => {
        api.removeRepo(getState().settings.currentUser, repo).then(
            () => {
                dispatch(showMessage("Success"));
                dispatch(removeRepoEntry(repo));
            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}

export function merge(theirs) {

    return (dispatch, getState) => {
        const settings = getSettings(getState());

        api.merge(settings, theirs).then(
            rs => {
                dispatch(showMessage("Success"));
            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}