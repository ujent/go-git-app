import * as api from './api';
import { ActionType, GetErrorMsg } from './constants';


export function getSettings(state) {
    return { user: state.settings.currentUser, repo: state.settings.currentRepo, branch: state.settings.currentBranch }
}

export function showSpinner() {
    return {
        type: ActionType.SHOW_SPINNER
    }
}

export function hideSpinner() {
    return {
        type: ActionType.HIDE_SPINNER
    }
}

export function showConfirm(msg, onConfirm, onClose) {
    return {
        type: ActionType.SHOW_CONFIRM,
        msg,
        onConfirm,
        onClose
    }
}

export function closeConfirm() {
    return {
        type: ActionType.CLOSE_CONFIRM,
    }
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
    const msg = GetErrorMsg(err);

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
        dispatch(resetMessage());

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
        dispatch(resetMessage());

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
        dispatch(resetMessage());

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

                if (rs.current) {
                    dispatch(getFiles(user, repo, rs.current));
                }
            },
            err => {
                dispatch(showError(err));
            }
        );
    }
}

export function switchBranch(name, isCreation) {

    return (dispatch, getState) => {
        const settings = getSettings(getState())

        if (isCreation) {
            if (!settings.branch) {
                dispatch(showMessage('Please, checkout branch!'))

                return;
            }
        }

        dispatch(showSpinner());
        dispatch(resetMessage());


        api.checkoutBranch(settings.user, settings.repo, name).then(
            () => {
                dispatch(setCurrentBranch(name));
                dispatch(getFiles(settings.user, settings.repo, name));
                dispatch(showMessage('Success'));
            },
            err => {
                dispatch(showError(err));
            }
        ).finally(
            () => dispatch(hideSpinner())
        );
    }
}

export function getFiles(user, repo, branch) {
    return (dispatch, getState) => {

        api.getFiles(user, repo, branch).then(
            rs => dispatch(setFiles(rs.files)),
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

export function getFile(path, isConflict, fileStatus) {

    return (dispatch, getState) => {
        dispatch(resetMessage());

        const settings = getSettings(getState())
        api.getFile(settings, path).then(
            rs => dispatch(setCurrentFile(rs.path, rs.content, isConflict, fileStatus)),
            err => {
                dispatch(showError(err));
            }
        );
    }
}

export function addFile(path, content) {
    return (dispatch, getState) => {
        dispatch(resetMessage());

        const state = getState();
        let hasFile = false;

        for (let i = 0; i < state.files.length; i++) {
            const f = state.files[i];

            if (f.path === path) {
                hasFile = true;
                break;
            }
        }

        if (hasFile) {
            dispatch(showError({ message: `File ${path} has already existed` }))
            return;
        }

        const settings = getSettings(state)

        dispatch(addFileEntry(path, content))

        api.addFile(settings, path, content).then(
            rs => {
            },
            err => {
                dispatch(removeFileEntry(path))
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

export function removeFile(path) {

    return (dispatch, getState) => {
        dispatch(resetMessage());

        const settings = getSettings(getState())
        api.removeFile(settings, path).then(
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


export function saveFile(path, content) {

    return (dispatch, getState) => {
        dispatch(showSpinner());
        dispatch(resetMessage());

        const settings = getSettings(getState())
        api.editFile(settings, path, content).then(
            rs => dispatch(showMessage('File was successfully saved')),
            err => {
                dispatch(showError(err));
            }
        ).finally(
            () => dispatch(hideSpinner())
        );
    }
}

export function setCurrentFile(name, content, isConflict, fileStatus) {
    return {
        type: ActionType.SET_CURRENT_FILE,
        name,
        content,
        isConflict,
        fileStatus
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
        dispatch(showSpinner());
        dispatch(resetMessage());

        const user = getState().settings.currentUser;

        api.clone(user, url, authName, authPsw).then(
            (rs) => {
                dispatch(switchRepo(rs.name));
                dispatch(showMessage('Success'));
            },
            err => {
                dispatch(showError(err));
            }
        ).finally(
            () => dispatch(hideSpinner())
        );

    }
}
export function commit(msg) {

    return (dispatch, getState) => {
        dispatch(showSpinner());
        dispatch(resetMessage());

        const settings = getSettings(getState())

        api.commit(settings, msg).then(

            () => {
                if (settings.branch) {
                    dispatch(getFiles(settings.user, settings.repo, settings.branch));
                } else {
                    dispatch(getBranches(settings.user, settings.repo));
                }
                dispatch(showMessage('Success'));
            },
            err => {
                dispatch(showError(err));
            }
        ).finally(
            () => dispatch(hideSpinner())
        );

    }
}
export function log() {

    return (dispatch, getState) => {
        dispatch(showSpinner());
        dispatch(resetMessage());

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
        ).finally(
            () => dispatch(hideSpinner())
        );

    }
}

export function pull(remote, authName, authPsw) {

    return (dispatch, getState) => {
        dispatch(showSpinner());
        dispatch(resetMessage());

        const settings = getSettings(getState())

        api.pull(settings, remote, authName, authPsw).then(
            (rs) => {
                const msg = rs.msg ? rs.msg : 'Success';

                dispatch(showMessage(msg));
                dispatch(getFiles(settings.user, settings.repo, settings.branch))
            },
            err => {
                dispatch(showError(err));
            }
        ).finally(
            () => dispatch(hideSpinner())
        );
    }
}
export function push(remote, authName, authPsw) {

    return (dispatch, getState) => {
        dispatch(showSpinner());
        dispatch(resetMessage());

        const settings = getSettings(getState())

        api.push(settings, remote, authName, authPsw).then(
            () => {
                dispatch(showMessage('Success'));
            },
            err => {
                dispatch(showError(err));
            }
        ).finally(
            () => dispatch(hideSpinner())
        );
    }
}
export function removeBranch(branch) {

    return (dispatch, getState) => {
        dispatch(resetMessage());

        const settings = getSettings(getState())

        api.removeBranch(settings.user, settings.repo, branch).then(
            () => {
                dispatch(showMessage('Success'));
                dispatch(removeBranchEntry(branch));
            },
            err => {
                dispatch(showError(err));
            }
        )
    }
}
export function removeRepo(repo) {

    return (dispatch, getState) => {
        dispatch(resetMessage());

        api.removeRepo(getState().settings.currentUser, repo).then(
            () => {
                dispatch(showMessage('Success'));
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
        dispatch(showSpinner());
        dispatch(resetMessage());

        const settings = getSettings(getState());

        api.merge(settings, theirs).then(
            rs => {

                if (rs.isFF) {
                    dispatch(showMessage('Fast-forward'));
                } else {
                    dispatch(showMessage(rs.msg))
                }

                dispatch(getFiles(settings.user, settings.repo, settings.branch))
            },
            err => {
                dispatch(showError(err));
            }
        ).finally(
            () => dispatch(hideSpinner())
        );

    }
}

export function abortMerge() {

    return (dispatch, getState) => {
        dispatch(showSpinner());
        dispatch(resetMessage());

        const settings = getSettings(getState());

        api.abortMerge(settings).then(
            rs => {
                dispatch(showMessage('Success'));
                dispatch(getFiles(settings.user, settings.repo, settings.branch))
            },
            err => {
                dispatch(showError(err));
            }
        ).finally(
            () => dispatch(hideSpinner())
        );
    }
}