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
                    dispatch(getFiles(user, repo, rs.currentBranch));
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
                dispatch(getFiles(settings.user, settings.repo, name));
                dispatch(showMessage("Success"));
            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}

export function getFiles(user, repo, branch) {
    return (dispatch, getState) => {

        api.getFiles(user, repo, branch).then(
            files => dispatch(setCurrentFile(files)),
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

export function getFile(path, isConflict) {
    return (dispatch, getState) => {

        const settings = getSettings(getState())
        api.getFile(settings, path, isConflict).then(
            rs => dispatch(setCurrentFile(rs.path, rs.content, isConflict)),
            err => {
                dispatch(showError(err));
            }
        );
    }
}

export function addFile(path, content) {
    return (dispatch, getState) => {

        const settings = getSettings(getState())
        api.addFile(settings, path, content).then(
            rs => dispatch(addFileEntry(path, content)),
            err => {
                dispatch(showError(err));
            }
        );
    }
}

export function addFileEntry(path, content) {
    return {
        type: ActionType.ADD_FILE_ENTRY,
        path,
        content
    }
}

export function removeFile(path, isConflict) {
    return (dispatch, getState) => {

        const settings = getSettings(getState())
        api.removeFile(settings, path, isConflict).then(
            rs => dispatch(removeFileEntry(path)),
            err => {
                dispatch(showError(err));
            }
        );
    }
}

export function removeFileEntry(path) {
    return {
        type: ActionType.REMOVE_FILE_ENTRY,
        path
    }
}


export function saveFile(path, content, isConflict) {
    return (dispatch, getState) => {

        const settings = getSettings(getState())
        api.editFile(settings, path, content, isConflict).then(
            rs => dispatch(showMessage("File was successfully saved")),
            err => {
                dispatch(showError(err));
            }
        );
    }
}

export function setCurrentFile(name, content, isConflict) {
    return {
        type: ActionType.SET_CURRENT_FILE,
        name,
        content,
        isConflict
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
                let res = ''
                rs.commits.forEach(function (el) {
                    res = res + `commit ${el.hash}
Author: ${el.author ? el.author.name : ''} ${el.author ? el.author.email : ''}
Date: ${el.date}
Msg: ${el.msg}
\n`
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