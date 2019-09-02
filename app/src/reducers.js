import { ActionType, StorageItem } from './constants';

export const rootReducer = (state = {}, action) => {
    switch (action.type) {
        case ActionType.SET_CURRENT_USER: {
            const prev = state.settings.currentUser;
            if (prev === action.user) {
                return state;
            }

            window.localStorage.setItem(StorageItem.User, action.user);

            return Object.assign({}, state, {
                repositories: [],
                branches: [],
                settings: Object.assign({}, state.settings, {
                    currentUser: action.user,
                    currentRepo: '',
                    currentBranch: ''
                }),
                files: [],
                currentFile: null
            });
        }
        case ActionType.SET_FILES: {
            return Object.assign({}, state, {
                files: action.files,
                currentFile: null
            })
        }
        case ActionType.SET_CURRENT_FILE: {
            let file = null;
            if (action.name) {
                file = {
                    path: action.name,
                    content: action.content,
                    isConflict: action.isConflict
                }
            }
            return Object.assign({}, state, {
                currentFile: file
            })
        }
        case ActionType.ADD_FILE_ENTRY: {
            const files = state.files;

            const f = {
                path: action.path,
                content: action.content,
                isConflict: false
            }

            files.unshift(f);

            return Object.assign({}, state, {
                currentFile: f,
                files: files
            });

        }
        case ActionType.REMOVE_FILE_ENTRY: {
            let current = state.currentFile;
            const files = state.files.filter(e => e.path !== action.path)

            if (current && current.path === action.path) {
                current = null;
            }

            return Object.assign({}, state, {
                files: files,
                currentFile: current
            });

        }
        case ActionType.SET_BRANCHES: {
            if (action.current) {
                window.localStorage.setItem(StorageItem.Branch, action.current);

                return Object.assign({}, state, {
                    branches: action.branches,
                    settings: Object.assign({}, state.settings, {
                        currentBranch: action.current,
                    })
                });
            }

            return Object.assign({}, state, {
                branches: action.branches
            });
        }
        case ActionType.SET_CURRENT_BRANCH: {
            if (state.settings.currentBranch === action.current) {
                return state;
            }

            window.localStorage.setItem(StorageItem.Branch, action.current);

            const branches = state.branches;

            if (branches.indexOf(action.current) !== -1) {
                return Object.assign({}, state, {
                    settings: Object.assign({}, state.settings, {
                        currentBranch: action.current,
                    })
                });
            }

            branches.push(action.current);

            return Object.assign({}, state, {
                branches: branches,
                settings: Object.assign({}, state.settings, {
                    currentBranch: action.current,
                })
            });
        }
        case ActionType.REMOVE_BRANCH: {
            const currentBranch = state.settings.currentBranch;
            let branches = state.branches;

            for (let i = 0; i < branches.length; i++) {
                if (branches[i] === action.branch) {
                    branches.splice(i, 1);
                    break;
                }
            }

            if (currentBranch === action.branch) {
                window.localStorage.removeItem(StorageItem.Branch);

                return Object.assign({}, state, {
                    branches: branches,
                    settings: Object.assign({}, state.settings, {
                        currentBranch: ''
                    }),
                    files: [],
                    currentFile: null
                });
            }

            return Object.assign({}, state, {
                branches: branches
            });
        }
        case ActionType.REMOVE_REPO: {
            const currentRepo = state.settings.currentRepo;
            let repos = state.repositories;

            for (let i = 0; i < repos.length; i++) {
                if (repos[i] === action.repo) {
                    repos.splice(i, 1);
                    break;
                }
            }

            if (currentRepo === action.repo) {
                window.localStorage.removeItem(StorageItem.Repo);
                window.localStorage.removeItem(StorageItem.Branch);

                return Object.assign({}, state, {
                    repositories: repos,
                    branches: [],
                    settings: Object.assign({}, state.settings, {
                        currentRepo: '',
                        currentBranch: ''
                    }),
                    currentFile: null,
                    files: []
                })
            }

            return Object.assign({}, state, {
                repositories: repos,
            })
        }
        case ActionType.SET_REPOSITORIES: {
            return Object.assign({}, state, {
                repositories: action.repos,
            });
        }
        case ActionType.SET_CURRENT_REPOSITORY: {
            if (state.settings.currentRepo === action.current) {
                return state;
            }

            window.localStorage.setItem(StorageItem.Repo, action.current);

            const repos = state.repositories;

            if (repos.indexOf(action.current) !== -1) {
                return Object.assign({}, state, {
                    settings: Object.assign({}, state.settings, {
                        currentRepo: action.current,
                    })
                });
            }

            repos.push(action.current);

            return Object.assign({}, state, {
                repositories: repos,
                settings: Object.assign({}, state.settings, {
                    currentRepo: action.current,
                })
            });

        }
        case ActionType.SHOW_MSG: {
            return Object.assign({}, state, {
                output: action.msg
            })
        }
        case ActionType.RESET_MSG: {
            return Object.assign({}, state, {
                output: ''
            })
        }
        default:
            return state;
    }
};