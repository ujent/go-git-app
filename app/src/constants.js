export const ActionType = {
    SHOW_MSG: 'SHOW_MSG',
    RESET_MSG: 'RESET_MSG',
    SHOW_SPINNER: 'SHOW_SPINNER',
    HIDE_SPINNER: 'HIDE_SPINNER',
    SHOW_CONFIRM: 'SHOW_CONFIRM',
    CLOSE_CONFIRM: 'CLOSE_CONFIRM',
    SET_CURRENT_USER: 'SET_CURRENT_USER',
    SET_REPOSITORIES: 'SET_REPOSITORIES',
    SET_CURRENT_REPOSITORY: 'SET_CURRENT_REPOSITORY',
    REMOVE_REPO: 'REMOVE_REPO',
    SET_CURRENT_BRANCH: 'SET_CURRENT_BRANCH',
    SET_BRANCHES: 'SET_BRANCHES',
    REMOVE_BRANCH: 'REMOVE_BRANCH',
    SET_FILES: 'SET_FILES',
    SET_CURRENT_FILE: 'SET_CURRENT_FILE',
    REMOVE_FILE_ENTRY: 'REMOVE_FILE_ENTRY',
    ADD_FILE_ENTRY: 'ADD_FILE_ENTRY'
}

export const StorageItem = {
    User: 'go_git_user',
    Repo: 'go_git_repo',
    Branch: 'go_git_branch'
}


export const FileStatus = {
    Unspecified: 0,
    Unmodified: 1,
    Modified: 2,
    Added: 3,
    Deleted: 4,
    Untracked: 5,
    Renamed: 6,
    Copied: 7,
    UpdatedButUnmerged: 8
}

export const GetErrorMsg = err => {
    let msg = '';

    if (err.status) {
        msg = `Error code: ${err.status}
Message: ${err.message}`
    } else {
        msg = err.message
    }

    return msg;
}